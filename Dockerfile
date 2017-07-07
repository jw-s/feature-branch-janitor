FROM golang
ADD . /go/src/github.com/JoelW-S/feature-branch-janitor
ADD ./vendor /go/src/github.com/JoelW-S/feature-branch-janitor/vendor
RUN go install github.com/JoelW-S/feature-branch-janitor/cmd/janitor
ENTRYPOINT ["janitor"]
