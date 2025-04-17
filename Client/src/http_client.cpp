#include "../include/http_client.h"
#include <iostream>

std::string fetchJsonFromUrl(const std::wstring& url) {
    URL_COMPONENTS urlComp = {0};
    urlComp.dwStructSize = sizeof(urlComp);
    urlComp.dwSchemeLength = (DWORD)-1;
    urlComp.dwHostNameLength = (DWORD)-1;
    urlComp.dwUrlPathLength = (DWORD)-1;

    if (!WinHttpCrackUrl(url.c_str(), (DWORD)url.length(), 0, &urlComp)) {
        std::cerr << "Failed to parse URL." << std::endl;
        return "";
    }

    std::wstring hostname(urlComp.lpszHostName, urlComp.dwHostNameLength);
    std::wstring path(urlComp.lpszUrlPath, urlComp.dwUrlPathLength);

    HINTERNET hSession = WinHttpOpen(
        L"WinHTTP Example/1.0", 
        WINHTTP_ACCESS_TYPE_DEFAULT_PROXY,
        WINHTTP_NO_PROXY_NAME, 
        WINHTTP_NO_PROXY_BYPASS, 
        0
    );
    if (!hSession) {
        std::cerr << "Failed to initialize WinHTTP." << std::endl;
        return "";
    }

    HINTERNET hConnect = WinHttpConnect(hSession, hostname.c_str(), urlComp.nPort, 0);
    if (!hConnect) {
        WinHttpCloseHandle(hSession);
        std::cerr << "Failed to connect to host." << std::endl;
        return "";
    }

    HINTERNET hRequest = WinHttpOpenRequest(
        hConnect, 
        L"GET", 
        path.empty() ? L"/" : path.c_str(),
        NULL, 
        WINHTTP_NO_REFERER, 
        WINHTTP_DEFAULT_ACCEPT_TYPES,
        (urlComp.nScheme == INTERNET_SCHEME_HTTPS) ? WINHTTP_FLAG_SECURE : 0
    );
    if (!hRequest) {
        WinHttpCloseHandle(hConnect);
        WinHttpCloseHandle(hSession);
        std::cerr << "Failed to create HTTP request." << std::endl;
        return "";
    }

    if (!WinHttpSendRequest(hRequest, NULL, 0, NULL, 0, 0, 0)) {
        WinHttpCloseHandle(hRequest);
        WinHttpCloseHandle(hConnect);
        WinHttpCloseHandle(hSession);
        std::cerr << "Failed to send HTTP request." << std::endl;
        return "";
    }

    if (!WinHttpReceiveResponse(hRequest, NULL)) {
        WinHttpCloseHandle(hRequest);
        WinHttpCloseHandle(hConnect);
        WinHttpCloseHandle(hSession);
        std::cerr << "Failed to receive HTTP response." << std::endl;
        return "";
    }

    std::string response;
    DWORD dwSize = 0;
    do {
        WinHttpQueryDataAvailable(hRequest, &dwSize);
        if (!dwSize) break;

        char* buffer = new char[dwSize + 1];
        DWORD dwDownloaded = 0;
        if (WinHttpReadData(hRequest, buffer, dwSize, &dwDownloaded)) {
            buffer[dwDownloaded] = '\0';
            response += buffer;
        }
        delete[] buffer;
    } while (dwSize > 0);

    WinHttpCloseHandle(hRequest);
    WinHttpCloseHandle(hConnect);
    WinHttpCloseHandle(hSession);

    return response;
}