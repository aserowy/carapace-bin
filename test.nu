let-env config = {
  show_banner: false

  edit_mode: vi
}

let carapace_completer = {|spans| 
  carapace $spans.0 nushell $spans | from json
}

mut current = (($env | default {} config).config | default {} completions)
$current.completions = ($current.completions | default {} external)
$current.completions.external = ($current.completions.external 
    | default true enable
    | default $carapace_completer completer)

let-env config = $current
