package protocol

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/iplist"
	humanize "github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

const clearScreen = "\033[H\033[2J"

const torrentBlockListURL = "http://john.bitsurge.net/public/biglist.p2p.gz"

var isHTTP = regexp.MustCompile(`^https?:\/\/`)

// ClientError formats errors coming from the client.
type ClientError struct {
	Type   string
	Origin error
}

func (clientError ClientError) Error() string {
	return fmt.Sprintf("Error %s: %s\n", clientError.Type, clientError.Origin)
}

// ClientConfig specifies the behaviour of a client.
type ClientConfig struct {
	TorrentPath    string
	Port           int
	TorrentPort    int
	Seed           bool
	TCP            bool
	MaxConnections int
}

// NewClientConfig creates a new default configuration.
func NewClientConfig() ClientConfig {
	return ClientConfig{
		Port:           8080,
		TorrentPort:    50007,
		Seed:           false,
		TCP:            true,
		MaxConnections: 200,
	}
}

// Client manages the torrent downloading.
type Client struct {
	Client   *torrent.Client
	Torrent  *torrent.Torrent
	Progress int64
	Uploaded int64
	Config   ClientConfig
	handle   SeekableContent
}

// NewClient creates a new torrent client based on a magnet or a torrent file.
// If the torrent file is on http, we try downloading it.
func NewClient(cfg ClientConfig) (client *Client, err error) {
	defer fmt.Printf("NewClient returned, torrent client created\n")
	var t *torrent.Torrent
	var c *torrent.Client
	client = &Client{}
	client.Config = cfg
	client.Client = c

	blocklist, err := loadBlocklist()
	if err != nil {
		return client, err
	}

	fmt.Printf("blocklist downloaded\n")
	torrentConfig := torrent.NewDefaultClientConfig()
	torrentConfig.DataDir = os.TempDir()
	torrentConfig.NoUpload = !cfg.Seed
	torrentConfig.Seed = cfg.Seed
	torrentConfig.DisableTCP = !cfg.TCP
	torrentConfig.ListenPort = cfg.TorrentPort
	torrentConfig.IPBlocklist = blocklist
	c, err = torrent.NewClient(torrentConfig)
	if err != nil {
		return client, ClientError{Type: "creating torrent client", Origin: err}
	}

	// Add as magnet url.
	if strings.HasPrefix(cfg.TorrentPath, "magnet:") {
		if t, err = c.AddMagnet(cfg.TorrentPath); err != nil {
			return client, ClientError{Type: "adding torrent", Origin: err}
		}
	} else {
		// Otherwise add as a torrent file.
		// If it's online, we try downloading the file.
		if isHTTP.MatchString(cfg.TorrentPath) {
			if cfg.TorrentPath, err = downloadFile(cfg.TorrentPath); err != nil {
				return client, ClientError{Type: "downloading torrent file", Origin: err}
			}
		}
		if t, err = c.AddTorrentFromFile(cfg.TorrentPath); err != nil {
			return client, ClientError{Type: "adding torrent to the client", Origin: err}
		}
	}
	client.Torrent = t
	client.Torrent.SetMaxEstablishedConns(cfg.MaxConnections)
	go func() {
		<-t.GotInfo()
		t.DownloadAll()
		// FIXME: How to priorize first 5% of the largest file, not the whole torrent.
		target := client.getLargestFile()
		target.SetPriority(torrent.PiecePriorityHigh)
		target.Torrent().DownloadPieces(0, int(t.NumPieces()/100*5))
		entry, err := NewFileReader(target)
		if err != nil {
			panic(errors.Wrap(err, "creating SeekableContent handle"))
		}
		client.handle = entry
	}()
	return client, err
}

// Close cleans up the connections.
func (c *Client) Close() (err error) {
	c.Torrent.Drop()
	c.Client.Close()
	if c.handle != nil {
		if e := c.handle.Close(); e != nil {
			err = e
		}
	}
	return nil
}

func (c *Client) Read(p []byte) (int, error) {
	if c.handle == nil {
		return 0, fmt.Errorf("resource not ready to read")
	}
	return c.handle.Read(p)
}

// Seek the resource.
func (c *Client) Seek(offset int64, whence int) (int64, error) {
	if c.handle == nil {
		return 0, fmt.Errorf("resource not ready to seek")
	}
	return c.handle.Seek(offset, whence)
}

// Status populates the status object with torrent stats.
func (c *Client) Status(s *Status) {
	if s == nil {
		return
	}
	if c.Torrent.Info() == nil {
		return
	}

	t := c.Torrent
	currentProgress := t.BytesCompleted()
	downloadSpeed := humanize.Bytes(uint64(currentProgress-c.Progress)) + "/s"
	stats := t.Stats()
	uploadProgress := (&stats).BytesWrittenData.Int64() - c.Uploaded

	s.Progress = currentProgress
	s.Uploaded = uploadProgress
	s.Throughput = downloadSpeed
	s.Size = t.Info().TotalLength()

	c.Progress = currentProgress
	c.Uploaded = uploadProgress
}

// ReadyForPlayback checks if the torrent is ready for playback or not.
// We wait until 5% of the torrent to start playing.
func (c Client) ReadyForPlayback() bool {
	return c.percentage() > 5
}

// Render outputs the command line interface for the client.
func (c *Client) Render() {
	t := c.Torrent

	if t.Info() == nil {
		return
	}

	currentProgress := t.BytesCompleted()
	downloadSpeed := humanize.Bytes(uint64(currentProgress-c.Progress)) + "/s"
	c.Progress = currentProgress

	complete := humanize.Bytes(uint64(currentProgress))
	size := humanize.Bytes(uint64(t.Info().TotalLength()))
	stats := t.Stats()
	uploadProgress := (&stats).BytesWrittenData.Int64() - c.Uploaded
	uploadSpeed := humanize.Bytes(uint64(uploadProgress)) + "/s"
	c.Uploaded = uploadProgress

	print(clearScreen)
	fmt.Println(t.Info().Name)
	fmt.Println(strings.Repeat("=", len(t.Info().Name)))
	if c.ReadyForPlayback() {
		fmt.Printf("Stream: \thttp://localhost:%d\n", c.Config.Port)
	}
	if currentProgress > 0 {
		fmt.Printf("Progress: \t%s / %s  %.2f%%\n", complete, size, c.percentage())
	}
	if currentProgress < t.Info().TotalLength() {
		fmt.Printf("Download speed: %s\n", downloadSpeed)
	}
	if c.Config.Seed {
		fmt.Printf("Upload speed: \t%s\n", uploadSpeed)
	}
}

func (c Client) getLargestFile() *torrent.File {
	var target *torrent.File
	var maxSize int64
	for _, file := range c.Torrent.Files() {
		if maxSize < file.Length() {
			maxSize = file.Length()
			target = file
		}
	}
	return target
}

func (c Client) percentage() float64 {
	info := c.Torrent.Info()
	if info == nil {
		return 0
	}
	return float64(c.Torrent.BytesCompleted()) / float64(info.TotalLength()) * 100
}

func downloadFile(URL string) (fileName string, err error) {
	var file *os.File
	if file, err = ioutil.TempFile(os.TempDir(), "novelty"); err != nil {
		return
	}
	defer func() {
		if ferr := file.Close(); ferr != nil {
			log.Printf("Error closing torrent file: %s", ferr)
		}
	}()
	response, err := http.Get(URL)
	if err != nil {
		return
	}
	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			log.Printf("Error closing torrent file: %s", ferr)
		}
	}()
	_, err = io.Copy(file, response.Body)
	return file.Name(), err
}

// Download and add the blocklist.
func loadBlocklist() (iplist.Ranger, error) {
	blocklistPath := filepath.Join(os.TempDir(), "novelty-blocklist.gz")
	if _, err := os.Stat(blocklistPath); os.IsNotExist(err) {
		if err := downloadBlockList(blocklistPath); err != nil {
			return nil, errors.Wrap(err, "downloading blocklist")
		}
	}
	blocklistReader, err := os.Open(blocklistPath)
	if err != nil {
		return nil, errors.Wrap(err, "opening blocklist")
	}
	gzipReader, err := gzip.NewReader(blocklistReader)
	if err != nil {
		return nil, errors.Wrap(err, "extracting blocklist")
	}
	blocklist, err := iplist.NewFromReader(gzipReader)
	if err != nil {
		return nil, errors.Wrap(err, "reading blocklist")
	}
	return blocklist, nil
}

func downloadBlockList(blocklistPath string) (err error) {
	log.Printf("Downloading blocklist")
	fileName, err := downloadFile(torrentBlockListURL)
	if err != nil {
		log.Printf("Error downloading blocklist: %s\n", err)
		return
	}
	return os.Rename(fileName, blocklistPath)
}

// SeekableContent describes an io.ReadSeeker that can be closed as well.
type SeekableContent interface {
	io.ReadSeeker
	io.Closer
}

// FileEntry helps reading a torrent file.
type FileEntry struct {
	*torrent.File
	torrent.Reader
}

// Seek seeks to the correct file position, paying attention to the offset.
func (f FileEntry) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.Seek(offset+f.File.Offset(), whence)
}

// NewFileReader sets up a torrent file for streaming reading.
func NewFileReader(f *torrent.File) (SeekableContent, error) {
	torrent := f.Torrent()
	reader := torrent.NewReader()

	// We read ahead 1% of the file continuously.
	reader.SetReadahead(f.Length() / 100)
	reader.SetResponsive()
	_, err := reader.Seek(f.Offset(), os.SEEK_SET)

	return &FileEntry{
		File:   f,
		Reader: reader,
	}, err
}
