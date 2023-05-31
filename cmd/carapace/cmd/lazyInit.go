package cmd

//go:generate go run ../../generate/gen.go

import (
	"fmt"
	"os"
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/xdg"
)

func bash_lazy(completers []string) string {
	snippet := `%v

_carapace_lazy() {
  source <(carapace $1 bash)
   $"_$1_completion"
}
complete -F _carapace_lazy %v
`
	return fmt.Sprintf(snippet, pathSnippet("bash"), strings.Join(completers, " "))
}

func bash_ble_lazy(completers []string) string {
	snippet := `%v

_carapace_lazy() {
  source <(carapace $1 bash-ble)
   $"_$1_completion_ble"
}
complete -F _carapace_lazy %v
`
	return fmt.Sprintf(snippet, pathSnippet("bash-ble"), strings.Join(completers, " "))
}

func elvish_lazy(completers []string) string {
	snippet := `%v

put %v | each {|c|
    set edit:completion:arg-completer[$c] = {|@arg|
        set edit:completion:arg-completer[$c] = {|@arg| }
        eval (carapace $c elvish | slurp)
        $edit:completion:arg-completer[$c] $@arg
    }
}
`
	return fmt.Sprintf(snippet, pathSnippet("elvish"), strings.Join(completers, " "))
}

func pathSnippet(shell string) (snippet string) {
	configDir, err := xdg.UserConfigDir()
	if err != nil {
		panic(err.Error())
	}
	binDir := configDir + "/carapace/bin"

	switch shell {
	case "bash", "bash-ble", "oil", "zsh":
		snippet = fmt.Sprintf(`export PATH="%v%v$PATH"`, binDir, string(os.PathListSeparator))

	case "elvish":
		snippet = fmt.Sprintf(`set paths = ['%v' $@paths]`, binDir)

	case "fish":
		snippet = fmt.Sprintf(`fish_add_path '%v'`, binDir)

	case "nushell":
		fixedBinDir := strings.ReplaceAll(binDir, `\`, `\\`)
		snippet = fmt.Sprintf(`
if "Path" in $env {
    let-env Path = ($env.Path | split row (char esep) | append "%v")
}

if "PATH" in $env {
    let-env PATH = ($env.PATH | split row (char esep) | append "%v")
}`,
			fixedBinDir, fixedBinDir)

	case "powershell":
		snippet = fmt.Sprintf(`[Environment]::SetEnvironmentVariable("PATH", "%v" + [IO.Path]::PathSeparator + [Environment]::GetEnvironmentVariable("PATH"))`, binDir)

	case "xonsh":
		snippet = fmt.Sprintf(`__xonsh__.env['PATH'].insert(0, '%v')`, binDir)

	default:
		snippet = fmt.Sprintf("# error: unknown shell: %#v", shell)
	}

	for _, path := range strings.Split(os.Getenv("PATH"), string(os.PathListSeparator)) {
		if path == binDir {
			carapace.LOG.Printf("PATH already contains %#v\n", binDir)
			if shell != "nushell" {
				snippet = "# " + snippet
			}
			break
		}
	}
	return
}

func fish_lazy(completers []string) string {
	snippet := `%v

function _carapace_lazy
   complete -c $argv[1] -e
   carapace $argv[1] fish | source
   complete --do-complete=(commandline -cp)
end
%v
`
	complete := make([]string, len(completers))
	for index, completer := range completers {
		complete[index] = fmt.Sprintf(`complete -c '%v' -f -a '(_carapace_lazy %v)'`, completer, completer)
	}
	return fmt.Sprintf(snippet, pathSnippet("fish"), strings.Join(complete, "\n"))
}

func nushell_lazy(completers []string) string {
	snippet := `%v

let carapace_completer = {|spans| 
  carapace $spans.0 nushell $spans | from json
}

mut current = (($env | default {} config).config | default {} completions)
$current.completions = ($current.completions | default {} external)
$current.completions.external = ($current.completions.external 
    | default true enable
    | default $carapace_completer completer)

let-env config = $current
    `

	return fmt.Sprintf(snippet, pathSnippet("nushell"))
}

func oil_lazy(completers []string) string {
	snippet := `%v

_carapace_lazy() {
  source <(carapace $1 oil)
   $"_$1_completion"
}
complete -F _carapace_lazy %v
`
	return fmt.Sprintf(snippet, pathSnippet("oil"), strings.Join(completers, " "))
}

func powershell_lazy(completers []string) string {
	snippet := `%v

$_carapace_lazy = {
    param($wordToComplete, $commandAst, $cursorPosition)
    $completer = $commandAst.CommandElements[0].Value
    carapace $completer powershell | Out-String | Invoke-Expression
    & (Get-Item "Function:_${completer}_completer") $wordToComplete $commandAst $cursorPosition
}
%v
`
	complete := make([]string, len(completers))
	for index, completer := range completers {
		complete[index] = fmt.Sprintf(`Register-ArgumentCompleter -Native -CommandName '%v' -ScriptBlock $_carapace_lazy`, completer)
	}
	return fmt.Sprintf(snippet, pathSnippet("powershell"), strings.Join(complete, "\n"))
}

func tcsh_lazy(completers []string) string {
	// TODO hardcoded for now
	snippet := make([]string, len(completers))
	for index, c := range completers {
		snippet[index] = fmt.Sprintf("complete \"%v\" 'p@*@`echo \"$COMMAND_LINE'\"''\"'\" | xargs carapace %v tcsh `@@' ;", c, c)
	}
	return strings.Join(snippet, "\n")
}

func xonsh_lazy(completers []string) string {
	snippet := `from xonsh.completers._aliases import _add_one_completer
from xonsh.completers.tools import *
import os

%v

@contextual_completer
def _carapace_lazy(context):
    """carapace lazy"""
    if (context.command and
        context.command.arg_index > 0 and
        context.command.args[0].value in [%v]):
        XSH.completers = XSH.completers.copy()
        exec(compile(subprocess.run(['carapace', context.command.args[0].value, 'xonsh'], stdout=subprocess.PIPE).stdout.decode('utf-8'), "", "exec"))
        return XSH.completers[context.command.args[0].value](context)
`
	complete := make([]string, len(completers))
	for index, completer := range completers {
		complete[index] = fmt.Sprintf(`'%v'`, completer)
	}
	snippet += `_add_one_completer('carapace_lazy', _carapace_lazy, 'start')`
	return fmt.Sprintf(snippet, pathSnippet("xonsh"), strings.Join(complete, ", "))
}

func zsh_lazy(completers []string) string {
	snippet := `%v

function _carapace_lazy {
    source <(carapace $words[1] zsh)
}
compdef _carapace_lazy %v
`
	return fmt.Sprintf(snippet, pathSnippet("zsh"), strings.Join(completers, " "))
}
