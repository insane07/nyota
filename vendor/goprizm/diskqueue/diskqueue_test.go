package diskqueue

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDiskQueueInit(t *testing.T) {
	assertDiskQueue := func(t *testing.T, opts Options) {
		os.RemoveAll(opts.Dir)
		dq, err := New(opts)
		defer dq.Close()

		if err != nil {
			t.Fatalf("failed to initialize err: %+v", err)
		}

		if dq.Name != filepath.Base(opts.Dir) {
			t.Fatalf("invalid disk queue name: %s exp: %s", dq.Name, filepath.Base(opts.Dir))
		}

		if dq.rSegmentID != opts.RecID.SegID {
			t.Fatalf("invalid segment id: %d exp: %d", dq.rSegmentID, opts.RecID.SegID)
		}

		if dq.rOffset != opts.RecID.Offset {
			t.Fatalf("invalid segment id: %d exp: %d", dq.rOffset, opts.RecID.Offset)
		}

		size := segSize
		if opts.SegmentSize != 0 {
			size = opts.SegmentSize
		}

		if dq.SegmentSize != size {
			t.Fatalf("invalid segment size: %d exp: %d", dq.SegmentSize, size)
		}
	}

	t.Run("Default-Options", func(t *testing.T) {
		opts := Options{
			Dir: "/tmp/diskqueue/test",
		}
		assertDiskQueue(t, opts)
	})

	t.Run("Custom-Options", func(t *testing.T) {
		opts := Options{
			Dir:         "/tmp/diskqueue/test",
			SegmentSize: 1024 * 1024,
			RecID: RecordID{
				SegID:  1,
				Offset: 1024,
			},
		}
		assertDiskQueue(t, opts)
	})

}

func TestEnqueue(t *testing.T) {
	opts := Options{
		Dir:         "/tmp/diskqueue/enqueue",
		SegmentSize: 10,
	}

	os.RemoveAll(opts.Dir)
	dq, err := New(opts)
	if err != nil {
		t.Fatalf("failed to initalize err: %+v", err)
	}

	t.Run("First-Segment", func(t *testing.T) {
		rec := []byte("Hello")
		dq.Enqueue(rec)
		dq.Close()
		dq = nil

		fileName := filepath.Join(opts.Dir, segmentName(1))
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			t.Fatalf("file %s does not exist err: %+v", fileName, err)
		}

		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			t.Fatalf("failed to read %s contents err: %+v", fileName, err)
		}

		if len(data) != (recHdrSize + len(rec)) {
			t.Fatalf("data(%+v) len(%d) mismatch. expected: %d", data, len(data), (recHdrSize + len(rec)))
		}

		var recLen int32

		if err := binary.Read(bytes.NewBuffer(data[:4]), binary.BigEndian, &recLen); err != nil {
			t.Fatalf("failed to convert record header to int32 err: %+v", err)
		}

		if int32(len(rec)) != recLen {
			t.Fatalf("len: %d mismatch expected: %d", recLen, len(rec))
		}

		if "Hello" != string(data[4:]) {
			t.Fatalf("expected: %s actual: %s", "Hello", string(data[4:]))
		}
	})

	t.Run("First-Segment-Repeat", func(t *testing.T) {
		dq, err := New(opts)
		if err != nil {
			t.Fatalf("failed to initalize err: %+v", err)
		}

		rec := []byte("World")
		dq.Enqueue(rec)
		dq.Close()
		dq = nil

		fileName := filepath.Join(opts.Dir, segmentName(1))
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			t.Fatalf("file %s does not exist err: %+v", fileName, err)
		}

		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			t.Fatalf("failed to read %s contents err: %+v", fileName, err)
		}

		if len(data) != (recHdrSize+len(rec))*2 {
			t.Fatalf("data(%+v) len(%d) mismatch. expected: %d", data, len(data), (recHdrSize+len(rec))*2)
		}

		var recLen int32

		if err := binary.Read(bytes.NewBuffer(data[9:13]), binary.BigEndian, &recLen); err != nil {
			t.Fatalf("failed to convert record header to int32 err: %+v", err)
		}

		if int32(len(rec)) != recLen {
			t.Fatalf("len: %d mismatch expected: %d", recLen, len(rec))
		}

		if "World" != string(data[13:]) {
			t.Fatalf("expected: %s actual: %s", "World", string(data[4:]))
		}
	})

	t.Run("Second-Third-Segment", func(t *testing.T) {
		dq, err := New(opts)
		if err != nil {
			t.Fatalf("failed to initalize err: %+v", err)
		}

		rec := []byte("Loreum\r\nIpsium")
		dq.Enqueue(rec)
		rec = []byte("Aruba\r\nPrizm")
		dq.Enqueue(rec)
		dq.Close()
		dq = nil

		files, err := ioutil.ReadDir(opts.Dir)
		if err != nil {
			t.Fatalf("failed to stat files err: %+v", err)
		}

		if len(files) != 3 {
			t.Fatalf("expeted %d segment files actual %d", 3, len(files))
		}

		var fileName string
		for _, fInfo := range files {
			if fInfo.Name() == segmentName(3) {
				fileName = fInfo.Name()
			}
		}

		if fileName == "" {
			t.Fatalf("file %s not found", fileName)
		}

		data, err := ioutil.ReadFile(filepath.Join(opts.Dir, fileName))
		if err != nil {
			t.Fatalf("failed to read %s contents err: %+v", fileName, err)
		}

		if len(data) != (recHdrSize + len(rec)) {
			t.Fatalf("data(%+v) len(%d) mismatch. expected: %d", data, len(data), (recHdrSize + len(rec)))
		}

		var recLen int32

		if err := binary.Read(bytes.NewBuffer(data[:4]), binary.BigEndian, &recLen); err != nil {
			t.Fatalf("failed to convert record header to int32 err: %+v", err)
		}

		if int32(len(rec)) != recLen {
			t.Fatalf("len: %d mismatch expected: %d", recLen, len(rec))
		}

		if "Aruba\r\nPrizm" != string(data[4:]) {
			t.Fatalf("expected: %s actual: %s", "Aruba\r\nPrizm", string(data[4:]))
		}

		if "Aruba\r\nPrizm" != string(data[4:]) {
			t.Fatalf("expected: %s actual: %s", "Aruba\r\nPrizm", string(data[4:]))
		}
	})
}

func TestDequeue(t *testing.T) {
	opts := Options{
		Dir:         "/tmp/diskqueue/dequeue",
		SegmentSize: 5,
	}

	os.RemoveAll(opts.Dir)
	dq, err := New(opts)
	if err != nil {
		t.Fatalf("failed to initialize err: %+v", err)
	}

	var recs []Record
	recs = append(recs, Record{RecordID{1, 0}, []byte("Hello"), RecordID{2, 0}})
	recs = append(recs, Record{RecordID{2, 0}, []byte("World"), RecordID{3, 0}})
	recs = append(recs, Record{RecordID{3, 0}, []byte("Loreum"), RecordID{4, 0}})

	for _, rec := range recs {
		dq.Enqueue(rec.Data)
	}

	dq.Close()
	dq = nil

	assertRecs := func(t *testing.T, expRec, actRec Record) {
		if string(expRec.Data) != string(actRec.Data) {
			t.Fatalf("expected: %s found: %s", string(expRec.Data), string(actRec.Data))
		}

		if expRec.ID != actRec.ID {
			t.Fatalf("expected ID: %+v found ID: %+v", expRec.ID, actRec.ID)
		}

		if expRec.NextID != actRec.NextID {
			t.Fatalf("expected nextID: %+v found nextID: %+v", expRec.NextID, actRec.NextID)
		}

	}

	t.Run("AllRec", func(t *testing.T) {
		dq, err := New(opts)
		if err != nil {
			t.Fatalf("failed to initialize err: %+v", err)
		}
		defer dq.Close()

		for i := 0; i < len(recs); i++ {
			r, err := dq.Dequeue()
			if err != nil {
				t.Fatalf("failed to read %+v err: %+v", recs[i], err)
			}

			assertRecs(t, recs[i], r)
		}

		dq = nil
	})

	t.Run("FirstRec", func(t *testing.T) {
		opts.RecID = RecordID{1, 0}
		dq, err := New(opts)
		if err != nil {
			t.Fatalf("failed to initialize err: %+v", err)
		}
		defer dq.Close()

		r, err := dq.Dequeue()
		if err != nil {
			t.Fatalf("failed to read %+v err: %+v", recs[0], err)
		}

		assertRecs(t, recs[0], r)
	})

	t.Run("SecondRec", func(t *testing.T) {
		opts.RecID = RecordID{2, 0}
		dq, err := New(opts)
		if err != nil {
			t.Fatalf("failed to initialize err: %+v", err)
		}
		defer dq.Close()

		r, err := dq.Dequeue()
		if err != nil {
			t.Fatalf("failed to read %+v err: %+v", recs[1], err)
		}

		assertRecs(t, recs[1], r)

	})

	t.Run("RepeatSecondRec", func(t *testing.T) {
		opts.RecID = RecordID{2, 0}
		dq, err := New(opts)
		if err != nil {
			t.Fatalf("failed to initialize err: %+v", err)
		}
		defer dq.Close()

		r, err := dq.Dequeue()
		if err != nil {
			t.Fatalf("failed to read %+v err: %+v", recs[1], err)
		}

		assertRecs(t, recs[1], r)

		dq.Seek(opts.RecID)
		r, err = dq.Dequeue()
		if err != nil {
			t.Fatalf("failed to read %+v err: %+v", recs[1], err)
		}

		assertRecs(t, recs[1], r)
	})

	t.Run("InvalidSegment", func(t *testing.T) {
		opts.RecID = RecordID{4, 10}
		dq, err := New(opts)
		if err != nil {
			t.Fatalf("failed to initialize err: %+v", err)
		}
		defer dq.Close()

		r, err := dq.Dequeue()
		if !os.IsNotExist(err) {
			t.Fatalf("expected no segment found err")
		}

		assertRecs(t, Record{}, r)
	})

	t.Run("InvalidFutureSegment", func(t *testing.T) {
		opts.RecID = RecordID{4, 10}
		dq, err := New(opts)
		if err != nil {
			t.Fatalf("failed to initialize err: %+v", err)
		}
		defer dq.Close()

		rec := Record{RecordID{4, 0}, []byte("Hiking"), RecordID{5, 0}}
		if err := dq.Enqueue(rec.Data); err != nil {
			t.Fatalf("failed to enqueue err: %+v", err)
		}

		if _, err = dq.Dequeue(); err == nil {
			t.Fatalf("expected err")
		}
	})

	t.Run("FutureSegment", func(t *testing.T) {
		opts.RecID = RecordID{5, 9}
		opts.SegmentSize = 18
		dq, err := New(opts)
		if err != nil {
			t.Fatalf("failed to initialize err: %+v", err)
		}
		defer dq.Close()

		rec := Record{RecordID{5, 9}, []byte("Trail"), RecordID{6, 0}}
		dq.Enqueue(rec.Data) // will insert at segment 4
		dq.Enqueue(rec.Data) // will insert at segment 5 offset 0
		// will insert at segment 5 offset 9
		if err := dq.Enqueue(rec.Data); err != nil {
			t.Fatalf("failed to enqueue err: %+v", err)
		}

		r, err := dq.Dequeue()
		if err != nil {
			t.Fatalf("failed to read %+v err: %+v", rec, err)
		}

		assertRecs(t, rec, r)
	})
}

func TestClean(t *testing.T) {
	opts := Options{
		Dir:         "/tmp/diskqueue/clean",
		SegmentSize: 5,
	}

	os.RemoveAll(opts.Dir)
	dq, err := New(opts)
	if err != nil {
		t.Fatalf("failed to initialize err: %+v", err)
	}

	var recs [][]byte
	recs = append(recs, []byte("Hello"))
	recs = append(recs, []byte("World"))
	recs = append(recs, []byte("Loreum"))

	for _, rec := range recs {
		dq.Enqueue(rec)
	}

	t.Run("ReatinAll", func(t *testing.T) {
		delCount, err := dq.Clean(3, false)
		if err != nil {
			t.Fatalf("failed to reatin all files err: %+v", err)
		}

		if delCount != 0 {
			t.Fatalf("expected: 0 deleted actual: %d", delCount)
		}
	})

	t.Run("DeleteFirst", func(t *testing.T) {
		delCount, err := dq.Clean(2, false)
		if err != nil {
			t.Fatalf("failed to delete file err: %+v", err)
		}

		if delCount != 1 {
			t.Fatalf("expected: 1 deleted actual: %d", delCount)
		}
	})

	// will delete 2.log and fail at 3.log
	t.Run("DeleteAllWithoutForce", func(t *testing.T) {
		delCount, err := dq.Clean(0, false)
		if err == nil {
			t.Fatalf("expected error")
		}

		if delCount != 1 {
			t.Fatalf("expected: 1 deleted actual: %d", delCount)
		}
	})

	t.Run("DeleteAllWithForce", func(t *testing.T) {
		delCount, err := dq.Clean(0, true)
		if err != nil {
			t.Fatalf("failed to delete all files err: %+v", err)
		}

		if delCount != 1 {
			t.Fatalf("expected: 1 deleted actual: %d", delCount)
		}
	})

	t.Run("Reset", func(t *testing.T) {
		if _, err := dq.Clean(0, true); err != nil {
			t.Fatalf("failed to delete all files err: %+v", err)
		}

		for _, rec := range recs {
			dq.Enqueue(rec)
		}

		for _, rec := range recs {
			r, err := dq.Dequeue()
			if err != nil {
				t.Fatalf("failed to read %+v err: %+v", rec, err)
			}

			if string(r.Data) != string(rec) {
				t.Fatalf("expected: %+v actual: %+v", string(rec), string(r.Data))
			}
		}
	})
}
