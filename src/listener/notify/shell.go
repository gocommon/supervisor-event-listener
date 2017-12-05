package notify

import "event"
import "os/exec"
import "strconv"

type Shell struct{}

// Send 会传递 eventname  playload ,  例/bin/echo eventname ip ProcessName GroupName FromState Expected Pid
func (shell *Shell) Send(message event.Message) error {

	cmd := exec.Command(Conf.Shell.Command, message.Header.EventName, message.Payload.Ip, message.Payload.ProcessName, message.Payload.GroupName, message.Payload.FromState, strconv.Itoa(message.Payload.Expected), strconv.Itoa(message.Payload.Pid))

	return cmd.Run()

}
