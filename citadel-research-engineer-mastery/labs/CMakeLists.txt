cmake_minimum_required(VERSION 3.20)
project(citadel_re_mastery LANGUAGES CXX)
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

if(NOT CMAKE_BUILD_TYPE)
  set(CMAKE_BUILD_TYPE Release)
endif()

add_executable(cache_benchmark cache_benchmark.cpp)
add_executable(branch_benchmark branch_benchmark.cpp)
add_executable(false_sharing false_sharing.cpp)
add_executable(spsc_queue spsc_queue.cpp)
add_executable(order_book order_book.cpp)

foreach(t cache_benchmark branch_benchmark false_sharing spsc_queue order_book)
  target_compile_options(${t} PRIVATE -Wall -Wextra -Wpedantic)
endforeach()
