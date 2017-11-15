NAME 			:=	stractl
VERSION		:=	$(shell git describe --tags)
REVISION	:=	$(shell git rev-parse --verify HEAD)
SRCS			:=	$(shell find . -type f -name '*.go')
LDFLAGS		:=	-ldflags "-X github.com/takuhiro/stractl/cmd.Version=$(shell git describe --tags) -X github.com/takuhiro/stractl/cmd.Revision=$(shell git rev-parse --verify HEAD)"

$(NAME): $(SRCS)
	go build $(LDFLAGS) -o $(NAME)
