# Behance Scraper API Service

<a href="https://behance.net"><img src="https://cdn.freebiesupply.com/logos/large/2x/behance-3-logo-png-transparent.png" width="180" alt="Selenium"/></a>

This is a Golang API service that provides functionality for scraping Behance.net and searching for projects based on user input.

## Utilities
* Golang: 1.20
* Playwright: 1.32 

## Project Structure

```
├── cmd
├── config
├── internal
│   ├── controllers
│   ├── errors
│   ├── middleware
│   ├── models
│   │   ├── requests
│   │   └── responses
│   ├── repositories
│   │   └── scraper
│   ├── services
│   └── utilities
├── pkg
└── server

16 directories
```

## Features
* Search Behance.net for projects based on keywords provided by the user
* Retrieve project details, including project image, title, author, total likes, and total views
* Flexible API design that can be integrated with various front-end applications or other backend services

## Example

**Endpoint**:
```
/behance/v1/search?search=searchitem
```
**Response**:

```
{
    "status": 200,
    "message": "Successfully fetched search result",
    "data": [
        {
            "imageUrl": "https://example.com",
            "projectUrl": "https://example.com",
            "title": "Example",
            "author": "Example",
            "totalLikes": 1000,
            "totalViews": 1000
        },
        {
            "imageUrl": "https://example.com",
            "projectUrl": "https://example.com",
            "title": "Example",
            "author": "Example",
            "totalLikes": 1000,
            "totalViews": 1000
        }
}
```

\* *Disclaimer:
This service is a scraper that retrieves information from Behance.net. I do not own Behance.net, nor affiliated with it in any way. By using this service, you agree to use it responsibly and in accordance with all applicable laws and regulations.*
