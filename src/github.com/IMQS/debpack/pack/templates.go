package pack

var templates map[string]string = map[string]string{
	"systemd": `[Unit]
Description=Job that runs the {{ .Name }} service
Documentation=man:{{ .Binary }}

[Service]
ExecStart=/usr/local/bin/{{ .Binary }}

[Install]
WantedBy=multi-user.target
`,
	"postinst": `#!/bin/sh -e
systemctl enable {{ .Binary }}.service
systemctl start {{ .Binary }}.service
`,
	"prerm": `#!/bin/sh -e
systemctl stop {{ .Binary }}.service
systemctl disable {{ .Binary }}.service`,
	"control": `Package: {{ .Binary }}
Version: {{ .Version }}
Section: {{ .Control.Section }}
Priority: {{ .Control.Priority }}
Architecture: {{ .Control.Architecture }}
Depends: {{ .Control.Depends }}
Maintainer: IMQS <imqs@imqs.co.za>
Description: {{ .Control.Description }}
{{ .Control.JoinedDescription }}
`,
}
