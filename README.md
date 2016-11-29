# prime

      Usage
      =====
      Web Service to consume a webpage, process some data and present it as JSON output

      Build - this will only work if you have go installed. (Tested on go version go1.7rc1 darwin/amd64)
      =====

      Clone repository.

      git clone git@github.com:nerfmiester/sainstest.git

      cd sainstest

      Choose your OS (It has only been tested on MAC)

      for MAC
      GOOS="darwin" go build sains.go

      for linux (64 bit)
      GOOS="linux" GOARCH="amd64" go build -o linuxsains sains.go

      for windows
      GOOS="windows" go build -o winsains.exe sains.go


      To run the tests

      go test -v -cover

      Execution
      ============

      Either if you have go installed

      go run sains.go

      Or execute binary
      _________________________

      linux
      Check permissions
      If the execution bit is not set

      chmod +x linuxsains

      ./linuxsains
      _________________________

      Mac

      ./sains
      _________________________

      windows

      winsains.exe
