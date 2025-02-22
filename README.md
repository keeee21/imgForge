# Image Upload & Resizer Web App

This is a simple web application that allows users to upload an image, resize it to a custom width and height, and download the resized image. The frontend is built with HTML, JavaScript, and the backend is powered by Go.

## 🚀 Features

- Upload an image (JPEG/PNG)
- Resize the image to a custom width and height
- Download the resized image
- Supports user-defined filename and save location (Chrome, Edge only)

## 🛠️ Technologies Used

- **Frontend**: HTML, JavaScript
- **Backend**: Go, Gorilla Mux, nfnt/resize, rs/cors
- **APIs**:
  - `fetch()` for image upload
  - `showSaveFilePicker()` for user-defined downloads (Chrome, Edge)

## 🚀 How to Run the Application
```sh
// The server will start at: http://localhost:8080
go run main.go
```