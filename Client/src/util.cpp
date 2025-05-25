#include "../include/util.h"
#include "../include/process_info.h"
#include <iostream>
#include <sys/types.h>
#include <signal.h>
#include <vector>
#include <string>
#include <Windows.h>
#include <TlHelp32.h>
#include <unordered_map>
#include <algorithm>
#include <chrono>


BYTE *buffer_payload(wchar_t *filename, OUT size_t &r_size)
{
    HANDLE file = CreateFileW(filename, GENERIC_READ, FILE_SHARE_READ, 0, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, 0);
    if(file == INVALID_HANDLE_VALUE) {
#ifdef _DEBUG
        std::cerr << "Could not open file!" << std::endl;
#endif
        return nullptr;
    }
    HANDLE mapping = CreateFileMapping(file, 0, PAGE_READONLY, 0, 0, 0);
    if (!mapping) {
#ifdef _DEBUG
        std::cerr << "Could not create mapping!" << std::endl;
#endif
        CloseHandle(file);
        return nullptr;
    }
    BYTE *dllRawData = (BYTE*) MapViewOfFile(mapping, FILE_MAP_READ, 0, 0, 0);
    if (dllRawData == nullptr) {
#ifdef _DEBUG
        std::cerr << "Could not map view of file" << std::endl;
#endif
        CloseHandle(mapping);
        CloseHandle(file);
        return nullptr;
    }
    r_size = GetFileSize(file, 0);
    BYTE* localCopyAddress = (BYTE*) VirtualAlloc(NULL, r_size, MEM_COMMIT | MEM_RESERVE, PAGE_READWRITE);
    if (localCopyAddress == NULL) {
        std::cerr << "Could not allocate memory in the current process" << std::endl;
        return nullptr;
    }
    memcpy(localCopyAddress, dllRawData, r_size);
    UnmapViewOfFile(dllRawData);
    CloseHandle(mapping);
    CloseHandle(file);
    return localCopyAddress;
}

void free_buffer(BYTE* buffer)
{
    if (buffer == NULL) return;
    VirtualFree(buffer, 0, MEM_RELEASE);
}

wchar_t* get_file_name(wchar_t *full_path)
{
    size_t len = wcslen(full_path);
    for (size_t i = len - 2; i >= 0; i--) {
        if (full_path[i] == '\\' || full_path[i] == '/') {
            return full_path + (i + 1);
        }
    }
    return full_path;
}


std::string GetWindowsUsername() {
    const DWORD MAX_USERNAME_LENGTH = 256;
    char username[MAX_USERNAME_LENGTH];
    DWORD size = MAX_USERNAME_LENGTH;
    
    if (!GetUserNameA(username, &size)) {
        std::cerr << "Error getting username. Code: " << GetLastError() << std::endl;
        return "unknown";
    }

    std::string result(username, size - 1);
    std::replace(result.begin(), result.end(), ' ', '-');
    
    return result;
}

wchar_t* get_directory(IN wchar_t *full_path, OUT wchar_t *out_buf, IN const size_t out_buf_size)
{
    memset(out_buf, 0, out_buf_size);
    memcpy(out_buf, full_path, out_buf_size);

    wchar_t *name_ptr = get_file_name(out_buf);
    if (name_ptr != nullptr) {
        *name_ptr = '\0'; //cut it
    }
    return out_buf;
}


bool IsPidRunning(DWORD pid) {
    HANDLE hProcess = OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, FALSE, pid);
    if (hProcess == NULL) {
        return false; // Process doesn't exist or access denied
    }

    DWORD exitCode;
    if (GetExitCodeProcess(hProcess, &exitCode)) {
        CloseHandle(hProcess);
        return (exitCode == STILL_ACTIVE); // Returns true if still running
    }

    CloseHandle(hProcess);
    return false;
}


bool AreProcessesRunning(const std::vector<std::string>& processNames) {
    HANDLE hProcessSnap;
    PROCESSENTRY32W pe32;
    
    // Take a snapshot of all processes in the system
    hProcessSnap = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
    if (hProcessSnap == INVALID_HANDLE_VALUE) {
        return false;
    }
    
    // Set the size of the structure before using it
    pe32.dwSize = sizeof(PROCESSENTRY32W);
    
    // Retrieve information about the first process
    if (!Process32FirstW(hProcessSnap, &pe32)) {
        CloseHandle(hProcessSnap); // clean the snapshot object
        return false;
    }
    
    // Setup converter from wide char to UTF-8
    std::wstring_convert<std::codecvt_utf8<wchar_t>, wchar_t> converter;
    
    // Walk through the process list
    do {
        // Convert the wide char process name to UTF-8 string
        std::string currentProcess = converter.to_bytes(pe32.szExeFile);
        std::transform(currentProcess.begin(), currentProcess.end(), currentProcess.begin(), ::tolower);
        
        // Check against each process in our list
        for (const auto& targetProcess : processNames) {
            std::string targetLower = targetProcess;
            std::transform(targetLower.begin(), targetLower.end(), targetLower.begin(), ::tolower);
            
            // If we find a match, return true
            if (currentProcess.find(targetLower) != std::string::npos) {
                CloseHandle(hProcessSnap);
                return true;
            }
        }
    } while (Process32NextW(hProcessSnap, &pe32));
    
    CloseHandle(hProcessSnap);
    return false;
}


// Convert std::string to LPWSTR
LPWSTR StringToLPWSTR(const std::string& str) {
    int size_needed = MultiByteToWideChar(CP_UTF8, 0, &str[0], (int)str.size(), NULL, 0);
    wchar_t* wstr = new wchar_t[size_needed + 1];
    MultiByteToWideChar(CP_UTF8, 0, &str[0], (int)str.size(), wstr, size_needed);
    wstr[size_needed] = 0;
    return wstr;
}

std::string buildCommandFromTemplate(
    const std::string& template_str,
    const std::unordered_map<std::string, std::string>& replacements
) {
    std::string result = template_str;
    
    for (const auto& [placeholder, value] : replacements) {
        size_t pos = result.find(placeholder);
        if (pos != std::string::npos) {
            result.replace(pos, placeholder.length(), value);
        }
    }
    
    return result;
}

bool IsDeviceIdle(int minutes) {
    // Get the last input time in milliseconds
    LASTINPUTINFO lastInputInfo;
    lastInputInfo.cbSize = sizeof(LASTINPUTINFO);
    
    if (!GetLastInputInfo(&lastInputInfo)) {
        std::cerr << "Failed to get last input info. Error: " << GetLastError() << std::endl;
        return false; // Assume not idle if we can't determine
    }

    // Calculate idle time in milliseconds
    DWORD currentTickCount = GetTickCount();
    DWORD idleTimeMs = currentTickCount - lastInputInfo.dwTime;

    // Convert minutes to milliseconds
    auto thresholdMs = std::chrono::minutes(minutes).count() * 60 * 1000;

    return (idleTimeMs >= thresholdMs);
}

bool IsForegroundWindowFullscreen() {
    HWND hwnd = GetForegroundWindow();
    if (!hwnd) return false;

    // Get window style to exclude certain windows
    LONG style = GetWindowLong(hwnd, GWL_STYLE);
    if (style & WS_CHILD) return false; // Ignore child windows

    RECT windowRect;
    GetWindowRect(hwnd, &windowRect);

    // Get the monitor where the window is located
    HMONITOR hMonitor = MonitorFromWindow(hwnd, MONITOR_DEFAULTTONEAREST);
    if (!hMonitor) return false;

    MONITORINFO monitorInfo;
    monitorInfo.cbSize = sizeof(monitorInfo);
    if (!GetMonitorInfo(hMonitor, &monitorInfo)) return false;

    // Check if window covers the entire monitor
    return (windowRect.left <= monitorInfo.rcMonitor.left &&
            windowRect.top <= monitorInfo.rcMonitor.top &&
            windowRect.right >= monitorInfo.rcMonitor.right &&
            windowRect.bottom >= monitorInfo.rcMonitor.bottom);
}



bool IsAnotherInstanceRunning(const char* mutexName) {
    HANDLE hMutex = CreateMutexA(
        NULL,           // Default security attributes
        TRUE,           // Initially owned
        mutexName);     // Unique mutex name (ANSI version)

    if (hMutex == NULL) {
        std::cerr << "CreateMutex error: " << GetLastError() << std::endl;
        return true; // Assume another instance is running to be safe
    }

    if (GetLastError() == ERROR_ALREADY_EXISTS) {
        ReleaseMutex(hMutex);
        CloseHandle(hMutex);
        return true;
    }

    // Mutex is held; release it when the program exits
    return false;
}