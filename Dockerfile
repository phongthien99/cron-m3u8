# First stage: build the executable.
FROM phongthien/golang-git:1.22 as build
ARG APP
# # It is important that these ARG's are defined after the FROM statement

# RUN echo "ARGNAME=${APP}"
# # git is required to fetch go dependencies
# RUN apt-get update && apt-get install git  -y


# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src
# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY apps/api/go.mod apps/api/go.mod ./
ENV GOSUMDB=off
RUN go mod download
# Import the code from the context.
COPY . .
# Build the executable to `/app`. Mark the build as statically linked.
WORKDIR /src/apps/api/cli
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /origin 

# RUN mv origin /src

# Final stage: the running container.
FROM gcr.io/distroless/base-debian10

WORKDIR /src

COPY --from=build /origin ./origin

EXPOSE 8080


CMD ["/src/origin","-c","/src/config/default.yaml"]
