ID{{"\t"}}Name{{"\t"}}Kind{{"\t"}}Archived{{"\t"}}Description
{{- range $channel := .}}
{{$channel.Id}}{{"\t"}}{{$channel.ChannelName}}{{"\t"}}{{$channel.Kind}}{{"\t"}}{{$channel.Archived}}{{"\t"}}{{$channel.Description}}
{{- end}}
