FROM golang:1.20 as build
WORKDIR /reservr
# Copy dependencies list
COPY . ./
RUN CGO_ENABLED=0 go build -tags lambda.norpc -o build/reservr-server main.go
# Copy artifacts to a clean image
FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /reservr/build/reservr-server ./reservr-server
ENTRYPOINT [ "./reservr-server" ]