About

You may have saved a file with the wrong extension, for example a JPEG file with a .psd extension. You cannot open the file in your editor, Or PSD save as jpeg extension, the saved file has a jpeg extension, but no software will open it.

If you cannot remember the actual file type the file should be you can use the following method to determine the true file type based on the file contents.

CGO_ENABLED=0 \
GOOS=linux \
GOARCH=amd64 \
GO111MODULE=off \
go build -ldflags "-X main.BUILD=`date -u +%F_%H:%M`_UTC" \
  -o search_filetype \
  main.go

./search_filetype -E -A -T -M ~/Documents
