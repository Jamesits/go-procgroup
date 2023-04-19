//go:build windows

package procgroup

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os/exec"
	"syscall"
	"unsafe"
)

type Group struct {
	hJob windows.Handle
}

func NewGroup() (ret *Group, err error) {
	ret = &Group{}
	ret.hJob, err = windows.CreateJobObject(nil, nil)

	return ret, err
}

func (g *Group) NewCmd() (*Cmd, error) {
	cmd := &Cmd{group: g}
	return cmd, nil
}

func (g *Group) Terminate(exitCode uint32) error {
	return windows.TerminateJobObject(g.hJob, exitCode)
}

type Cmd struct {
	exec.Cmd
	group *Group
}

func (c *Cmd) Start() (err error) {
	if c.group == nil {
		panic("you are expected to create the cmd from a group")
	}

	if c.Cmd.SysProcAttr == nil {
		c.Cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	originalCreationFlags := c.Cmd.SysProcAttr.CreationFlags
	// https://learn.microsoft.com/en-us/windows/win32/api/jobapi2/nf-jobapi2-assignprocesstojobobject
	// If the process is being monitored by the Program Compatibility Assistant (PCA), it is placed into a compatibility job. Therefore, the process must be created using CREATE_BREAKAWAY_FROM_JOB before it can be placed in another job.
	c.Cmd.SysProcAttr.CreationFlags |= windows.CREATE_SUSPENDED | windows.CREATE_BREAKAWAY_FROM_JOB

	err = c.Cmd.Start()
	if err != nil {
		return fmt.Errorf("start: %w", err)
	}

	// put the process into a job

	hProcess, err := windows.OpenProcess(windows.PROCESS_SET_QUOTA|windows.PROCESS_TERMINATE, false, uint32(c.Cmd.Process.Pid))
	if err != nil {
		terminateProcess(uint32(c.Cmd.Process.Pid))
		return fmt.Errorf("OpenProcess: %w", err)
	}

	err = windows.AssignProcessToJobObject(c.group.hJob, hProcess)
	if err != nil {
		terminateProcess(uint32(c.Cmd.Process.Pid))
		return fmt.Errorf("AssignProcessToJobObject: %w", err)
	}

	err = windows.CloseHandle(hProcess)
	if err != nil {
		terminateProcess(uint32(c.Cmd.Process.Pid))
		return fmt.Errorf("CloseHandle(process): %w", err)
	}

	// find the main thread of the new process, and resume it

	if (originalCreationFlags & windows.CREATE_SUSPENDED) != 0 {
		return nil
	}

	hSnapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPTHREAD, uint32(c.Cmd.Process.Pid))
	if err != nil {
		terminateProcess(uint32(c.Cmd.Process.Pid))
		return fmt.Errorf("CreateToolhelp32Snapshot: %w", err)
	}

	// https://learn.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-createtoolhelp32snapshot
	threadEntry := &windows.ThreadEntry32{}
	threadEntry.Size = uint32(unsafe.Sizeof(*threadEntry))
	err = windows.Thread32First(hSnapshot, threadEntry)
	if err != nil {
		terminateProcess(uint32(c.Cmd.Process.Pid))
		return fmt.Errorf("Thread32First: %w", err)
	}
	for {
		if threadEntry.OwnerProcessID == uint32(c.Cmd.Process.Pid) {
			break
		}

		threadEntry.Size = uint32(unsafe.Sizeof(*threadEntry))
		err = windows.Thread32Next(hSnapshot, threadEntry)
		if err != nil {
			terminateProcess(uint32(c.Cmd.Process.Pid))
			return fmt.Errorf("Thread32Next: %w", err)
		}
	}

	err = windows.CloseHandle(hSnapshot)
	if err != nil {
		terminateProcess(uint32(c.Cmd.Process.Pid))
		return fmt.Errorf("CloseHandle(snapshot): %w", err)
	}

	hThread, err := windows.OpenThread(windows.THREAD_SUSPEND_RESUME, false, threadEntry.ThreadID)
	if err != nil {
		terminateProcess(uint32(c.Cmd.Process.Pid))
		return fmt.Errorf("OpenThread: %w", err)
	}

	_, err = windows.ResumeThread(hThread)
	if err != nil {
		terminateProcess(uint32(c.Cmd.Process.Pid))
		return fmt.Errorf("ResumeThread: %w", err)
	}

	err = windows.CloseHandle(hThread)
	if err != nil {
		terminateProcess(uint32(c.Cmd.Process.Pid))
		return fmt.Errorf("CloseHandle(thread): %w", err)
	}

	return nil
}

func terminateProcess(pid uint32) {
	hProcess, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, pid)
	if err != nil {
		return
	}
	defer windows.CloseHandle(hProcess)

	_ = windows.TerminateProcess(hProcess, 1)
}
