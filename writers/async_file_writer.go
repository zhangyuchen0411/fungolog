package writers

type AsyncFileWriter struct {
	w     *FileWriter
	c     chan []byte
	close chan bool
}

func NewAsyncFileWriter(w *FileWriter, bufSize int) *AsyncFileWriter {
	aw := &AsyncFileWriter{
		w:     w,
		c:     make(chan []byte, bufSize),
		close: make(chan bool),
	}
	go aw.asyncWrite()
	return aw
}

func (this *AsyncFileWriter) Write(data []byte) error {
	this.c <- data
	return nil
}

func (this *AsyncFileWriter) Close() error {
	this.close <- true
	return nil
}

func (this *AsyncFileWriter) asyncWrite() {
	var data []byte
	for {
		select {
		case <-this.close:
			this.w.Close()
			break
		case data = <-this.c:
			this.w.Write(data)
		}
	}
}

var _ Writer = &AsyncFileWriter{}