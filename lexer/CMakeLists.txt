cmake_minimum_required(VERSION 4.0)
project(lexer)

set(CMAKE_CXX_STANDARD 20)

add_library(lexer
        src/Lexer.cpp
)

target_include_directories(lexer PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/include)

add_executable(Lexer src/main.cpp
        src/Lexer.cpp)
target_link_libraries(Lexer PRIVATE lexer)

