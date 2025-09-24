cmake_minimum_required(VERSION 3.20)
project(OpenSSLQUICBugReproducer C)

set(CMAKE_C_STANDARD 11)

set(OpenSSL_DIR  ${CMAKE_SOURCE_DIR}/3rd/OpenSSL/lib64/cmake/OpenSSL)
find_package(OpenSSL CONFIG REQUIRED)
get_target_property(OpenSSL_INCLUDE_DIR OpenSSL::SSL INTERFACE_INCLUDE_DIRECTORIES)
include_directories(${OpenSSL_INCLUDE_DIR})
message(STATUS ${OpenSSL_INCLUDE_DIR})

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -fstack-usage -fsanitize=address,undefined -fno-omit-frame-pointer -g -D_FORTIFY_SOURCE=2")
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fstack-usage -fsanitize=address,undefined -fno-omit-frame-pointer -g -D_FORTIFY_SOURCE=2")

add_executable(client
        src/main.c)
target_link_libraries(client PRIVATE
        OpenSSL::SSL
        OpenSSL::Crypto
        pthread
        )
