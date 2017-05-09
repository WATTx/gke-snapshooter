Google Compute Disks Snapshot Report:

{{ range $i, $status := . }}
{{ $status.Disk }} - {{ $status.Text }}
{{ end }}
