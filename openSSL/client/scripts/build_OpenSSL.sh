#!/bin/bash
set -eu
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
readonly openssl_version="3.5.3"
readonly my_3rd_library_path="$HOME/Documents/dev/c/3rd"
readonly working_dir="/tmp/openssl_building"
readonly install_path="${script_dir}/../3rd/OpenSSL"
readonly build_thread=$(nproc)
readonly clone_OpenSSL=false
readonly checkout_OpenSSL=false
readonly common_configure_options=(
    "--api=3.0"
    "-static"
    "no-acvp-tests"
    "no-afalgeng"
    "no-apps"
    "no-argon2"
    "no-aria"
    "no-asan"
    # "asm"
    "no-async"
    "no-atexit"
    "no-autoalginit"
    "no-autoerrinit"
    "no-autoload-config"
    "no-bf"
    "no-blake2"
    "no-brotli"
    "no-brotli-dynamic"
    "no-buildtest-c++"
    # "bulk"
    "no-cached-fetch"
    "no-camellia"
    "no-capieng"
    "no-winstore"
    "no-cast"
    # "chacha"
    "no-cmac"
    "no-cmp"
    "no-cms"
    "no-comp"
    "no-crypto-mdebug"
    "no-ct"
    # "default-thread-pool"
    "no-demos"
    "no-h3demo"
    "no-hqinterop"
    "no-deprecated"
    "no-des"
    "no-devcryptoeng"
    # "dgram"
    "no-dh"
    "no-docs"
    "no-dsa"
    "no-dso"
    # "dtls"
    "no-dynamic-engine"
    # "ec"
    "no-ec2m"
    "no-ec_nistp_64_gcc_128"
    # "ecdh"
    # "ecdsa"
    # "ecx"
    "no-egd"
    "no-engine"
    "no-err"
    "no-external-tests"
    "no-filenames"
    "no-fips"
    "no-fips-securitychecks"
    "no-fips-post"
    "no-fips-jitter"
    "no-fuzz-afl"
    "no-fuzz-libfuzzer"
    "no-gost"
    "no-http"
    "no-idea"
    "no-integrity-only-ciphers"
    "no-jitter"
    "no-ktls"
    "no-legacy"
    "no-loadereng"
    "no-makedepend"
    "no-md2"
    "no-md4"
    "no-mdc2"
    "no-ml-dsa"
    "no-ml-kem"
    "no-module"
    "no-msan"
    "no-multiblock"
    "no-nextprotoneg"
    "no-ocb"
    "no-ocsp"
    "no-padlockeng"
    # "pic"
    # "pie"
    "no-pinshared"
    # "poly1305"
    # "posix-io"
    "no-psk"
    # "quic"
    "no-unstable-qlog"
    "no-rc2"
    "no-rc4"
    "no-rc5"
    "no-rdrand"
    "no-rfc3779"
    "no-rmd160"
    "no-scrypt"
    "no-sctp"
    "no-secure-memory"
    "no-seed"
    "no-shared"
    "no-siphash"
    "no-siv"
    "no-slh-dsa"
    "no-sm2"
    "no-sm2-precomp"
    "no-sm3"
    "no-sm4"
    # "sock"
    "no-srp"
    # "srtp"
    "no-sse2"
    # "ssl"
    "no-ssl-trace"
    "no-static-engine"
    # "stdio"
    "no-sslkeylog"
    "no-tests"
    "no-tfo"
    # "thread-pool"
    # "threads"
    # "tls"
    "no-tls1"
    "no-tls1_1"
    "no-tls1_1-method"
    # "tls1_2"
    # "tls1_2-method"
    "no-tls-deprecated-ec"
    "no-trace"
    "no-ts"
    "no-ubsan"
    "no-ui-console"
    "no-unit-test"
    "no-uplink"
    "no-weak-ssl-ciphers"
    "no-whirlpool"
    "no-zlib"
    "no-zlib-dynamic"
    "no-zstd"
    "no-zstd-dynamic"
)

print_info() {
    echo -e "\033[32m[INFO]\033[0m $1"
}

print_warning() {
    echo -e "\033[33m[WARNING]\033[0m $1"
}

print_error() {
    echo -e "\033[31m[ERROR]\033[0m $1"
}

git_clone_OpenSSL(){
    rm -fr "${my_3rd_library_path}/openssl_src"
    mkdir -p "${my_3rd_library_path}"
    git clone https://github.com/openssl/openssl.git "${my_3rd_library_path}/openssl_src"
}

git_checkout_OpenSSL(){
    cd "$my_3rd_library_path/openssl_src" || exit
    git checkout -f "openssl-$openssl_version"
}

clean_and_create_dirs(){
    rm -fr "${install_path}"
    mkdir -p "${install_path}"

    rm -fr "${working_dir}"
    mkdir -p "${working_dir}"
}

print_path_info(){
    print_info "OpenSSL source path: ${my_3rd_library_path}/openssl_src"
    print_info "OpenSSL build path: ${working_dir}"
    print_info "OpenSSL install path: ${install_path}"
}
build_OpenSSL(){
    if [ "$clone_OpenSSL" = true ]; then
        git_clone_OpenSSL
    fi
    if [ "$checkout_OpenSSL" = true ]; then
        git_checkout_OpenSSL
    fi

    clean_and_create_dirs
    cd "$working_dir" || exit
    print_info "prepare to build OpenSSL, output directory: ${install_path}"
    export CFLAGS="-O0 -g3 -ggdb3 -fno-omit-frame-pointer -DDEBUG -UNDEBUG -fPIC"
    export CXXFLAGS="-O0 -g3 -ggdb3 -fno-omit-frame-pointer -DDEBUG -UNDEBUG -fPIC"
    configure_options=(
        "--prefix=${install_path}"
        "linux-x86_64"
        "${common_configure_options[@]}"
    )

    "$my_3rd_library_path/openssl_src/Configure" "${configure_options[@]}"
    make "-j${build_thread}"
    make install
}
print_path_info
# build_OpenSSL