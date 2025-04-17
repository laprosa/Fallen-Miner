#pragma once
#include <string>
#include <windows.h>
#include <winhttp.h>

std::string fetchJsonFromUrl(const std::wstring& url);