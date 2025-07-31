package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	url    string
	output string
)

var styleError = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF5555")).
	Bold(true)

var styleSuccess = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#50FA7B")).
	Bold(true)

var rootCmd = &cobra.Command{
	Use:   "imgdl",
	Short: "Downloads an image with a progress bar",
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			fmt.Fprintln(os.Stderr, styleError.Render("Error: --url is required"))
			cmd.Usage()
			os.Exit(1)
		}

		if err := downloadImageWithProgress(url, output); err != nil {
			fmt.Fprintln(os.Stderr, styleError.Render("Download failed: "+err.Error()))
			os.Exit(1)
		}

		fmt.Println(styleSuccess.Render("Downloaded successfully to " + output))
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL of the image to download (required)")
	rootCmd.Flags().StringVarP(&output, "output", "o", "output.jpg", "Output filename")
	_ = rootCmd.MarkFlagRequired("url")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, styleError.Render("Command failed: "+err.Error()))
		os.Exit(1)
	}
}

// --------------- Download Logic + Progress Integration --------------------

func downloadImageWithProgress(url, output string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64)")
	req.Header.Set("Referer", "https://uhdpaper.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	outFile, err := os.Create(output)
	if err != nil {
		resp.Body.Close()
		return err
	}
	defer outFile.Close()

	m := NewProgressModel(resp.ContentLength, outFile)
	m.respBody = resp.Body // assign the reader here!

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

// ---------------------- Bubble Tea Progress Model ------------------------

type ProgressMsg float64
type DownloadDoneMsg struct{}
type DownloadErrorMsg struct{ Err error }

type progressModel struct {
	progress progress.Model
	percent  float64
	total    int64
	writer   io.Writer
	done     bool
	url      string
	respBody io.ReadCloser
}

func NewProgressModel(total int64, writer io.Writer) progressModel {
	return progressModel{
		progress: progress.New(
			progress.WithGradient("#00ffcc", "#0066ff"),
			progress.WithWidth(50),
		),
		total:  total,
		writer: writer,
	}
}

func (m progressModel) Init() tea.Cmd {
	return m.download()
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProgressMsg:
		m.percent = float64(msg)
		// Schedule next chunk read command:
		return m, m.readChunk(m.respBody, m.writer, m.total, int64(m.percent*float64(m.total)))
	case DownloadDoneMsg:
		m.done = true
		return m, tea.Quit
	case DownloadErrorMsg:
		fmt.Fprintln(os.Stderr, styleError.Render("Error: "+msg.Err.Error()))
		return m, tea.Quit
	}
	return m, nil
}

func (m progressModel) View() string {
	if m.done {
		return styleSuccess.Render("\n✅ Download complete\n")
	}
	bar := m.progress.ViewAs(m.percent)
	percent := fmt.Sprintf("%.0f%%", m.percent*100)
	return lipgloss.NewStyle().Padding(1).Render(fmt.Sprintf("Downloading... %s\n%s", percent, bar))
}

func (m progressModel) download() tea.Cmd {
	return m.readChunk(m.respBody, m.writer, m.total, 0)
}

func (m progressModel) readChunk(r io.Reader, w io.Writer, total int64, written int64) tea.Cmd {
	return func() tea.Msg {
		buf := make([]byte, 32*1024)
		n, err := r.Read(buf)
		if n > 0 {
			if _, werr := w.Write(buf[:n]); werr != nil {
				return DownloadErrorMsg{Err: werr}
			}
			written += int64(n)
			progress := float64(written) / float64(total)
			if progress > 1 {
				progress = 1
			}
			return ProgressMsg(progress)
		}
		if err == io.EOF {
			return DownloadDoneMsg{}
		}
		if err != nil {
			return DownloadErrorMsg{Err: err}
		}
		return ProgressMsg(float64(written) / float64(total))
	}
}

func readAndTrackProgress(r io.Reader, w io.Writer, total int64) tea.Msg {
	buf := make([]byte, 32*1024)
	var written int64

	for {
		n, err := r.Read(buf)
		if n > 0 {
			if _, werr := w.Write(buf[:n]); werr != nil {
				return DownloadErrorMsg{Err: werr}
			}
			written += int64(n)
			progress := float64(written) / float64(total)
			if progress > 1 {
				progress = 1
			}
			if total > 0 {
				// yield message
				return ProgressMsg(progress)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return DownloadErrorMsg{Err: err}
		}
	}

	return DownloadDoneMsg{}
}
