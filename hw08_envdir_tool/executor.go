package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var args []string
	osEnv := os.Environ()

	switch len(cmd) {
	case 1:
		args = make([]string, 0)
	default:
		args = cmd[1:]
	}
	commandName := cmd[0]

	command := exec.Command(commandName, args...)
	command.Env = prepareEnvs(osEnv, env)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr

	err := command.Start()
	if err != nil {
		return command.ProcessState.ExitCode()
	}
	log.Printf("Waiting for command to finish...")
	err = command.Wait() // ошибка выполнения
	if err != nil {
		return command.ProcessState.ExitCode()
	}

	return 0
}

func prepareEnvs(osEnv []string, editions Environment) []string {
	envs := make(map[string]string)

	for i := range osEnv {
		row := strings.Split(osEnv[i], "=")
		envs[row[0]] = row[1]
	}

	for k, v := range editions {
		if v.NeedRemove {
			delete(envs, k)
			continue
		}
		envs[k] = v.Value
	}

	res := make([]string, 0)
	for k, v := range envs {
		res = append(res, k+"="+v)
	}

	return res
}
