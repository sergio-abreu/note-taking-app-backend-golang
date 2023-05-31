package cron

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"

	"github.com/gofrs/uuid"

	"github.com/sergio-abreu/note-taking-app-backend-golang/domain/notes"
)

func NewLocalCron(baseDir string) LocalCron {
	return LocalCron{baseDir: baseDir}
}

type LocalCron struct {
	baseDir string
}

func (l LocalCron) CreateCron(ctx context.Context, reminder notes.Reminder) error {
	err := l.createCronFile(reminder)
	if err != nil {
		return err
	}

	err = l.updateCrontab()
	if err != nil {
		return err
	}

	return nil
}

func (l LocalCron) DeleteCron(ctx context.Context, reminder notes.Reminder) error {
	if reminder.ID == uuid.Nil {
		return nil
	}

	err := l.removeCronFile(reminder)
	if err != nil {
		return err
	}

	err = l.updateCrontab()
	if err != nil {
		return err
	}

	return nil
}

func (l LocalCron) removeCronFile(reminder notes.Reminder) error {
	err := os.Remove(fmt.Sprintf("%s/%s.cron", l.baseDir, reminder.ID))
	if err != nil {
		return err
	}
	return nil
}

func (l LocalCron) createCronFile(reminder notes.Reminder) error {
	err := os.Mkdir(l.baseDir, os.ModeDir|os.ModePerm)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	f, err := os.Create(fmt.Sprintf("%s/%s.cron", l.baseDir, reminder.ID))
	if err != nil {
		return err
	}
	defer f.Close()
	cron := reminder.ParseCron()
	id := reminder.ID
	endsAt := reminder.ParseEndsAt(cron)
	// if [[ -z '123' ]] || [[ "$(date +%s)" -le "$(date -d 2013-07-18T01:00:00Z +%s)" ]]; then echo greater; else echo smaller; fi
	_, err = f.Write([]byte(
		fmt.Sprintf(`%s curl -x POST "$NOTE_TAKING_BASE_URL"/v1/webhook/reminder/%s?date=%s
`, cron, id, endsAt)))
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (l LocalCron) updateCrontab() error {
	dirEntries, err := os.ReadDir(l.baseDir)
	if err != nil {
		return err
	}
	var files []string
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			files = append(files, path.Join(l.baseDir, dirEntry.Name()))
		}
	}
	cmd := exec.Command("cat", files...)
	crontab, err := cmd.Output()
	if err != nil {
		return err
	}

	cmd = exec.Command("crontab", "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()

	_, err = stdin.Write(crontab)
	if err != nil {
		return err
	}
	stdin.Close()

	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
