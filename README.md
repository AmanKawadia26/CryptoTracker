# CryptoTracker

CryptoTracker is a command-line application built in Go that allows users to track cryptocurrency prices, set alerts, and manage user profiles. The application features an intuitive user interface with enhanced visuals using the Fatih color library.

✨ Features
User Authentication: Login and Signup functionalities.
Admin and User Panels: Separate panels for admin and regular users.
View Top Cryptocurrencies: Easily access the top 10 cryptocurrencies.
Search Functionality: Look up specific cryptocurrencies by name.
Price Alerts: Set custom price alerts for your favorite cryptocurrencies.
Unavailable Crypto Tracking: Request tracking for cryptocurrencies not currently supported.
⚙️ Installation
Prerequisites
Ensure that you have Go installed on your system. You can download it from the official Go website.

Clone the Repository - 
git clone https://github.com/AmanKawadia26/CryptoTracker.git
cd CryptoTracker

Retrieve Dependencies - 
Before running the project, ensure you have all the necessary Go packages by running the following command:
go get -u

Make sure to include the following packages:
Fatih Color Library: This library is used to enhance the user interface with colored output in the terminal.

go get -u github.com/fatih/color

⚠️ Important Notes on File Paths
There are specific functions in the project where you need to manually set the file path according to your device:

GetAllUsers() in request.go: This function reads user data from a file. Ensure you adjust the file path to match your device's directory structure.

LoadConfig() in config.go: This function loads the application's configuration from a file. Similar to GetAllUsers(), the path should be modified according to your file location.

Make sure to review and update these paths before running the application.

🚀 Running the Application
Once the dependencies are set up and file paths are adjusted, you can run the application with:

go run main.go

This command will compile and execute the main.go file, starting the CryptoTracker application.

📜 License
This project is licensed under the MIT License - see the LICENSE file for details.

