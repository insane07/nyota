// append only queue backed by disk as storage
//
// Enqueue and Dequeue APIs are not for concurrent access. Caller has to synchronize the access.
// Clean API provides option for caller to reatin max segments.
//
// Records
// =======
// record is stored as <len><data> format in segment file.
//
// Segement Files
// ==============
// Records are stored in fixed sized files called segments. New segments are created
// when size of current segments exceeds limit. Segments file names are of the form
// {segment-id}.log. It start with 1(1.log) and is incremented for every new segment.
// Caller has to reset the queue if segment file size is changed.
// If the  segment size is reduced after enqueuing the data we will miss reading records.
//
// RecordID
// =========
// Each record has a unique id which is of the form {segment-id}:{offset}.
// Records can be located using this unique id.

package diskqueue

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"goprizm/log"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var (
	segExt     = ".log"
	segFmt     = "%d" + segExt
	segSize    = 1024 * 1024 * 10 // default segment file size - 10MB
	recHdrSize = 4                // 4 bytes
)

type Options struct {
	Dir         string   // Dir to store segment file
	SegmentSize int      // Max size of segement file
	RecID       RecordID // Index of the record to start reading
}

type DiskQueue struct {
	Options
	Name string // Name of diskqueue obtained from dir name.

	wFile      *os.File     // write - currently opened file
	wSegmentID int          // write - id of current segment
	wSize      int          // write - size of current segment
	wBuf       bytes.Buffer // write - buffio for record

	rFile      *os.File // read - currently openet file
	rOffset    int      // read - offset of the current record
	rSegmentID int      // read - id of current segment
}

type RecordID struct {
	SegID  int `json:"segment_id"` // segment ID
	Offset int `json:"offset"`     // offset of the record in the segment
}

func (rID RecordID) String() string {
	return fmt.Sprintf("%d:%d", rID.SegID, rID.Offset)
}

type Record struct {
	ID     RecordID // record id
	Data   []byte   // record data
	NextID RecordID // id of the next record
}

func (r *Record) IsEmpty() bool {
	return len(r.Data) == 0
}

func (r Record) String() string {
	return fmt.Sprintf("ID: %s Data: %s NextID: %s", r.ID, string(r.Data), r.NextID)
}

type Stats struct {
	Name        string `json:"name"`
	Dir         string `json:"dir"`
	SegmentSize int    `json:"segment_size"`

	SegmentCount   int `json:"segments"`
	FirstSegmentID int `json:"first_segment_id"`
	LastSegmentID  int `json:"last_segment_id"`

	RSegmentID int `json:"read_segment_id"`
	ROffset    int `json:"read_offset"`

	WSegmentID   int `json:"write_segment_id"`
	WSegmentSize int `json:"write_segment_size"`
}

func New(opts Options) (*DiskQueue, error) {
	if err := os.MkdirAll(opts.Dir, 0755); err != nil {
		return nil, err
	}

	dq := &DiskQueue{
		Options:    opts,
		Name:       filepath.Base(opts.Dir),
		rSegmentID: opts.RecID.SegID,
		rOffset:    opts.RecID.Offset,
	}

	if dq.SegmentSize == 0 {
		dq.SegmentSize = segSize
	}

	log.Printf("diskqueue(%s) - init with segment size: %d", dq.Name, dq.SegmentSize)

	if err := dq.seekReadSegment(); err != nil {
		return nil, err
	}
	return dq, nil
}

func (dq *DiskQueue) seekReadSegment() error {
	// no RecordID mentioned to start
	if dq.rSegmentID == 0 {
		return nil
	}

	if err := dq.openReadSegment(); err != nil {
		if os.IsNotExist(err) {
			log.Printf("diskqueue(%s) - read segment: %s does not exist", dq.Name, segmentName(dq.rSegmentID))
			return nil
		}
		return err
	}

	if _, err := dq.rFile.Seek(int64(dq.rOffset), 0); err != nil {
		return err
	}

	log.Debugf("diskqueue(%s) - segment: %s seek: %d for read", dq.Name, dq.rFile.Name(), dq.rOffset)
	return nil
}

func (dq *DiskQueue) openReadSegment() error {
	var err error
	file := filepath.Join(dq.Dir, segmentName(dq.rSegmentID))
	dq.rFile, err = os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}

	if _, err := dq.rFile.Seek(int64(dq.rOffset), 0); err != nil {
		return err
	}

	log.Debugf("diskqueue(%s) - open: %s segment with offset: %d for read", dq.Name, file, dq.rOffset)

	return nil
}

// Enqueue record
func (dq *DiskQueue) Enqueue(rec []byte) error {
	if err := dq.prepareEnqueue(); err != nil {
		return err
	}

	size, err := dq.enqueue(rec)
	if err != nil {
		return err
	}

	//log.Debugf("diskqueue(%s) - enqueue rec: %+v", dq.Name, rec)
	dq.wSize += size

	return nil
}

func (dq *DiskQueue) prepareEnqueue() error {
	if dq.wSegmentID == 0 {
		// get edge segement ids
		_, lastSegID, err := dq.getEdgeSegmentIDs()
		if err != nil {
			return err
		}

		dq.wSegmentID = lastSegID
		// create first file
		if dq.wSegmentID == 0 {
			dq.wSegmentID = 1
			return dq.openWriteSegment()
		}

		stat, err := os.Stat(filepath.Join(dq.Dir, segmentName(dq.wSegmentID)))
		if err != nil {
			return err
		}

		dq.wSize = int(stat.Size())
		//create next file
		if dq.wSize >= dq.SegmentSize {
			dq.wSize = 0
			dq.wSegmentID += 1
			return dq.openWriteSegment()
		}

		// open last file
		return dq.openWriteSegment()
	}

	//create next file
	if dq.wSize >= dq.SegmentSize {
		if dq.wFile != nil {
			dq.wFile.Close()
			dq.wFile = nil
		}

		dq.wSize = 0
		dq.wSegmentID += 1
		return dq.openWriteSegment()
	}

	// handle permission issue
	if dq.wFile == nil {
		return dq.openWriteSegment()
	}

	// current file
	return nil
}

func (dq *DiskQueue) getEdgeSegmentIDs() (int, int, error) {
	segments, err := dq.getSegmentNames()
	if err != nil {
		return 0, 0, err
	}

	if len(segments) == 0 {
		return 0, 0, nil
	}

	first, err := segmentID(segments[0])
	if err != nil {
		return 0, 0, err
	}

	last, err := segmentID(segments[len(segments)-1])
	if err != nil {
		return 0, 0, err
	}

	return first, last, nil
}

func (dq *DiskQueue) getSegmentNames() ([]string, error) {
	files, err := ioutil.ReadDir(dq.Dir)
	if err != nil {
		return []string{}, err
	}

	var segments []string
	for _, file := range files {
		if filepath.Ext(file.Name()) != segExt {
			continue
		}

		segments = append(segments, file.Name())
	}

	sort.Slice(segments, func(i, j int) bool {
		f0, err := segmentID(segments[i])
		if err != nil {
			return false
		}

		f1, err := segmentID(segments[j])
		if err != nil {
			return false
		}

		return f0 < f1

	})

	return segments, nil
}

func (dq *DiskQueue) openWriteSegment() error {
	var err error
	file := filepath.Join(dq.Dir, segmentName(dq.wSegmentID))
	dq.wFile, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	log.Debugf("diskqueue(%s) - open: %s segment for write", dq.Name, file)

	return nil
}

func (dq *DiskQueue) enqueue(rec []byte) (int, error) {
	dq.wBuf.Reset()
	if err := binary.Write(&dq.wBuf, binary.BigEndian, int32(len(rec))); err != nil {
		return 0, err
	}

	if _, err := dq.wBuf.Write(rec); err != nil {
		return 0, err
	}

	if _, err := dq.wFile.Write(dq.wBuf.Bytes()); err != nil {
		return 0, err
	}

	return dq.wBuf.Len(), nil
}

// Dequeue record
func (dq *DiskQueue) Dequeue() (Record, error) {
	//no read file opened
	if dq.rFile == nil {

		// get first file from dir
		if dq.rSegmentID == 0 {
			firstSegID, _, err := dq.getEdgeSegmentIDs()
			if err != nil {
				return Record{}, err
			}

			dq.rSegmentID = firstSegID

			//no records enqueued
			if dq.rSegmentID == 0 {
				return Record{}, nil
			}
		}

		if err := dq.openReadSegment(); err != nil {
			if os.IsNotExist(err) {
				log.Printf("diskqueue(%s) - read segment: %s does not exist", dq.Name, segmentName(dq.rSegmentID))
				return Record{}, err
			}
			return Record{}, err
		}
	}

	rec, err := dq.dequeue()
	if err != nil {
		return Record{}, err
	}

	//log.Debugf("diskqueue(%s) - dequeue rec: %+v", dq.Name, rec)
	return rec, err
}

func (dq *DiskQueue) dequeue() (Record, error) {
	moveToNextSegment := func() {
		dq.rFile.Close()
		dq.rFile = nil

		dq.rOffset = 0
		dq.rSegmentID += 1
		log.Debugf("diskqueue(%s) - move to next read segment: %s", dq.Name, filepath.Join(dq.Dir, segmentName(dq.rSegmentID)))
	}

	// read record length
	var recLen int32
	if err := binary.Read(dq.rFile, binary.BigEndian, &recLen); err != nil {
		if err == io.EOF {
			_, last, err := dq.getEdgeSegmentIDs()
			if err != nil {
				return Record{}, err
			}

			// next segment exists and current file is EOF
			if last > dq.rSegmentID {
				moveToNextSegment()
			}
		}
		return Record{}, err
	}

	rec := make([]byte, recLen)
	if _, err := io.ReadFull(dq.rFile, rec); err != nil {
		/*
			if err == io.EOF {
				log.Debugf("diskqueue(%s) - read segment: %s EOF. Reset the queue", dq.Name, dq.rFile.Name())
			}
		*/
		return Record{}, err
	}

	// segment id and offset of the current record
	recSegID := dq.rSegmentID
	recOffset := dq.rOffset

	// offset of the next record
	dq.rOffset += (recHdrSize + int(recLen))

	// move to next segement
	// if the  segment size is reduced after enqueuing the data we will miss reading records
	if dq.rOffset >= dq.SegmentSize {
		moveToNextSegment()
	}

	recID := RecordID{SegID: recSegID, Offset: recOffset}
	nextRecID := RecordID{SegID: dq.rSegmentID, Offset: dq.rOffset}

	return Record{ID: recID, Data: rec, NextID: nextRecID}, nil
}

func (dq *DiskQueue) Seek(recID RecordID) error {
	dq.rSegmentID = recID.SegID
	dq.rOffset = recID.Offset
	return dq.seekReadSegment()
}

// Retain max number of segment files.
// If segement is open for reading or writing ignore unles force option is set
// Return count of segment deleted
func (dq *DiskQueue) Clean(retainMax int, force bool) (int, error) {
	segments, err := dq.getSegmentNames()
	if err != nil {
		return 0, err
	}

	// number of segments available less than max required
	if len(segments) <= retainMax {
		log.Debugf("diskqueue(%s) - segment count: %d less than max retain count: %d skip clean", dq.Name, len(segments), retainMax)
		return 0, nil
	}

	delList := segments[:len(segments)-retainMax]
	var delCount int
	for _, segment := range delList {
		segFile := filepath.Join(dq.Dir, segment)

		if dq.rFile != nil && dq.rFile.Name() == segFile {
			if !force {
				return delCount, fmt.Errorf("diskqueue(%s) - read segment: %s is opened", dq.Name, segFile)
			}

			dq.rFile.Close()
			dq.rFile = nil
			log.Printf("diskqueue(%s) - read segment: %s force closed", dq.Name, segFile)
		}

		if dq.wFile != nil && dq.wFile.Name() == segFile {
			if !force {
				return delCount, fmt.Errorf("diskqueue(%s) - write segment: %s is opened", dq.Name, segFile)
			}

			dq.wFile.Close()
			dq.wFile = nil
			dq.wSize = 0
			log.Printf("diskqueue(%s) - write segment: %s force closed", dq.Name, segFile)
		}

		if err := os.Remove(segFile); err != nil {
			return delCount, err
		}

		delCount += 1
	}

	// all the read and write segments are deleted. start fresh
	if retainMax == 0 && force {
		dq.resetSegmentCounter()
	}

	log.Debugf("diskqueue(%s) - segment max retain: %d total: %d delete: %d", dq.Name, retainMax, len(segments), delCount)
	return delCount, nil
}

// reset the read and write segment counter
func (dq *DiskQueue) resetSegmentCounter() {
	dq.rSegmentID = 0
	dq.rOffset = 0

	dq.wSegmentID = 0
	log.Debugf("diskqueue(%s) - reset segment counter", dq.Name)
}

func (dq *DiskQueue) Close() {
	if dq.wFile != nil {
		fileName := dq.wFile.Name()
		dq.wFile.Close()
		dq.wFile = nil
		log.Printf("diskqueue(%s) - close write segment: %s", dq.Name, fileName)
	}

	if dq.rFile != nil {
		fileName := dq.rFile.Name()
		dq.rFile.Close()
		dq.rFile = nil
		log.Printf("diskqueue(%s) - close read segment: %s", dq.Name, fileName)
	}
}

func segmentName(id int) string {
	return fmt.Sprintf(segFmt, id)
}

// segmentID extracts id int from file name.
func segmentID(file string) (int, error) {
	fields := strings.Split(file, ".")
	if len(fields) != 2 || fields[1] != segExt[1:] {
		return 0, fmt.Errorf("invalid fmt(ext-%s)", file)
	}

	return strconv.Atoi(fields[0])
}

func (dq *DiskQueue) Stats() (Stats, error) {
	var stats Stats

	stats.Name = dq.Name
	stats.Dir = dq.Dir
	stats.SegmentSize = dq.SegmentSize

	segNames, err := dq.getSegmentNames()
	if err != nil {
		return stats, err
	}

	stats.SegmentCount = len(segNames)

	firstSegID, lastSegID, err := dq.getEdgeSegmentIDs()
	if err != nil {
		return stats, err
	}

	stats.FirstSegmentID = firstSegID
	stats.LastSegmentID = lastSegID

	stats.RSegmentID = dq.rSegmentID
	stats.ROffset = dq.rOffset

	stats.WSegmentID = dq.wSegmentID
	stats.WSegmentSize = dq.wSize

	return stats, nil
}
