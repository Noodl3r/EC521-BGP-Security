cmake_minimum_required(VERSION 3.20)
set(TARGET_NAME project)

project(project VERSION 0.1.0 LANGUAGES C)

set(CMAKE_C_STANDARD 17)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_C_EXTENSIONS OFF)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Manually add source directories here.
# You otherwise will need to glob and those can be inconsistent.
add_executable(${TARGET_NAME} src/main.c)

target_include_directories(
    ${TARGET_NAME}
    PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_compile_options(${TARGET_NAME} PRIVATE -Wall -Wextra -Wpedantic -O3)

if(CMAKE_BUILD_TYPE STREQUAL "Debug")
    target_compile_options(${TARGET_NAME} PRIVATE -g3 -O0)
endif()
