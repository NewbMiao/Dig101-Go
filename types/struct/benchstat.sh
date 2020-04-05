#! /bin/env bash
set -e

# command -v go >/dev/null 2>&1 || {
#     echo "installing go..."
#     wget -O go.tar.gz -L -q https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
#     sudo tar -C /usr/local -xzf go.tar.gz
#     echo "export PATH=$PATH:/usr/local/go/bin:$(/usr/local/go/bin/go env GOPATH)/bin" > ~/.bashrc
#     export PATH=$PATH:/usr/local/go/bin:$(/usr/local/go/bin/go env GOPATH)/bin
#     rm -f go.tar.gz

#     command -v git >/dev/null 2>&1 || {
#         apt-get -qq update -y && apt-get -qq install -y git
#     }
#     echo  "installing benchstat..."
#     go get -u golang.org/x/perf/cmd/benchstat
# }


UsePerfLock=""

command -v benchstat >/dev/null 2>&1 || { 
    echo "installing benchstat..."
    go get -u golang.org/x/perf/cmd/benchstat  
}

if [ -f "/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq" ]; then
    command -v perflock >/dev/null 2>&1 || { 
        echo  "May support perflock, try use perflock"
        go get -u github.com/aclements/perflock/cmd/perflock 
        perflock -daemon & 
        UsePerfLock="perflock"
    }
fi

# sh benchstat.sh 20
N=${1:-1}
echo "Will run address-align bench $N times"

workspace=$(cd $(dirname $0) && pwd -P)
{
    cd $workspace

    GOGC=off GODEBUG=asyncpreemptoff=1 $UsePerfLock go test -gcflags='-N -l' -run none . -bench . -count $N -cpu 1 > b.txt
    benchstat b.txt
}
