cmake_minimum_required(VERSION 3.14)
project(vible_cpp_sdk LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

if(NOT CMAKE_BUILD_TYPE AND NOT CMAKE_CONFIGURATION_TYPES)
  set(CMAKE_BUILD_TYPE Release CACHE STRING "" FORCE)
endif()

# ── libcurl (system) ───────────────────────────────────────────────────────
# HTTP / SSE / audio upload / calls all go through libcurl. Required at build
# time — install it via your platform package manager if missing:
#   macOS : brew install curl       (or use the system libcurl)
#   Debian: apt-get install libcurl4-openssl-dev
find_package(CURL REQUIRED)

# ── nlohmann/json (header-only, fetched) ───────────────────────────────────
# Pulled via FetchContent so we don't vendor the giant single header into the
# repo. Needs network access on first configure.
include(FetchContent)
FetchContent_Declare(
  nlohmann_json
  URL https://github.com/nlohmann/json/releases/download/v3.11.3/json.tar.xz
  URL_HASH SHA256=d6c65aca6b1ed68e7a182f4757257b107ae403032760ed6ef121c9d55e81757d
  DOWNLOAD_EXTRACT_TIMESTAMP TRUE
)
FetchContent_MakeAvailable(nlohmann_json)

# Header-only "client" target: just the include dir + its dependencies.
add_library(vible INTERFACE)
target_include_directories(vible INTERFACE ${CMAKE_CURRENT_SOURCE_DIR}/include)
target_link_libraries(vible INTERFACE CURL::libcurl nlohmann_json::nlohmann_json)

# ── samples ────────────────────────────────────────────────────────────────
foreach(sample chat stream models speech transcribe converse call)
  add_executable(${sample} samples/${sample}.cpp)
  target_link_libraries(${sample} PRIVATE vible)
endforeach()
