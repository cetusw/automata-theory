cmake_minimum_required(VERSION 4.0)
project(RecursiveDescent)

set(CMAKE_CXX_STANDARD 20)

add_library(recursive-descent
        src/main.cpp
)

target_include_directories(recursive-descent PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/include)

add_executable(RecursiveDescent src/main.cpp)
target_link_libraries(RecursiveDescent PRIVATE recursive-descent)
