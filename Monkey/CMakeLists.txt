cmake_minimum_required(VERSION 4.0)
project(Monkey)

set(CMAKE_CXX_STANDARD 20)

add_library(monkey
        src/main.cpp
)

target_include_directories(monkey PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/include)

add_executable(Monkey src/main.cpp)
target_link_libraries(Monkey PRIVATE monkey)

