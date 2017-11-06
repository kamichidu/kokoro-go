ID: {{.Id}}
Name: {{.ChannelName}}
Kind: {{.Kind}}
Archived: {{.Archived}}
Description:
{{.Description}}

Members:
{{- range $membership := .Memberships}}
{{$membership.Authority}}{{"\t"}}{{$membership.Profile.Id}}{{"\t"}}{{$membership.Profile.Type}}{{"\t"}}{{$membership.Profile.ScreenName}}{{"\t"}}({{$membership.Profile.DisplayName}})
{{- end}}
