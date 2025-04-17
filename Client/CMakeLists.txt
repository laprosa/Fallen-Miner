cmake_minimum_required(VERSION 3.10)
project(fallen-miner)

# Force GCC (MinGW) if not already set
if(NOT CMAKE_C_COMPILER)
    set(CMAKE_C_COMPILER gcc)
endif()
if(NOT CMAKE_CXX_COMPILER)
    set(CMAKE_CXX_COMPILER g++)
endif()


set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)


add_definitions(-DUNICODE -D_UNICODE)


# Source files
set(SRC_FILES
    src/main.cpp
    src/util.cpp
    src/pe_hdrs_helper.cpp
    src/hollowing_parts.cpp
    src/delete_pending_file.cpp
    src/http_client.cpp
    src/json_printer.cpp
    src/process_info.cpp
)

# Header files
set(HEADER_FILES
    include/util.h
    include/kernel32_undoc.h
    include/ntddk.h
    include/pe_hdrs_helper.h
    include/hollowing_parts.h
    include/delete_pending_file.h
    include/http_client.h
    include/json_printer.h
    include/json.hpp
    include/process_info.h
    include/resources.h
    include/embedded_resource.h
)

# Resource files
set(RESOURCE_FILES
    resources/embedded_resource.rc
    resources/embedded.exe
)

# Create executable
add_executable(fallen-miner 
    WIN32 
    ${SRC_FILES} 
    ${HEADER_FILES} 
    ${RESOURCE_FILES}
)

target_compile_options(fallen-miner PRIVATE -w)

# Resource dependencies
set_source_files_properties(
    ${CMAKE_CURRENT_SOURCE_DIR}/resources/embedded_resource.rc
    PROPERTIES
    OBJECT_DEPENDS "${CMAKE_CURRENT_SOURCE_DIR}/resources/embedded.exe"
)

# Link libraries (GCC/MinGW style)
target_link_libraries(fallen-miner
    -lwinhttp
    -lntdll
)

# Include directories
target_include_directories(fallen-miner PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/include
)

# Installation
install(TARGETS fallen-miner
    DESTINATION ${CMAKE_INSTALL_PREFIX}
    COMPONENT fallen-miner
)

add_custom_command(TARGET fallen-miner POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E remove_directory ${CMAKE_BINARY_DIR}/CMakeFiles
    COMMAND ${CMAKE_COMMAND} -E remove ${CMAKE_BINARY_DIR}/CMakeCache.txt
    COMMAND ${CMAKE_COMMAND} -E remove ${CMAKE_BINARY_DIR}/cmake_install.cmake
    COMMAND ${CMAKE_COMMAND} -E remove ${CMAKE_BINARY_DIR}/Makefile
    COMMENT "Auto-removing CMake-generated files"
)

message(STATUS "Project configured for GCC/MinGW")