## Set Execution Policy to Unrestricted if needed
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy Unrestricted

Write-Host "Starting the script..."

# See current working directory
# Get-Location

# Change Directory to server
Set-Location .\server\

## Running of Authentication Service
Write-Host "Starting Authentication Service..."
Set-Location .\auth-service\
Start-Process -NoNewWindow -FilePath "powershell.exe" -ArgumentList "go run .\auth_service.go"
cd ..

## Running of Payment Service
Write-Host "Starting Payment Service..."
Set-Location .\payment-service\
Start-Process -NoNewWindow -FilePath "powershell.exe" -ArgumentList "go run .\payment_service.go"
cd ..

# ## Running of User Management Service
Set-Location .\user-management-service\
Write-Host "Starting User Management Service..."
Start-Process -NoNewWindow -FilePath "powershell.exe" -ArgumentList "go run .\user-management-service.go"
cd ..

## Running of Vehicle Service
Set-Location .\vehicle-service\
Write-Host "Starting Vehicle Service..."
Start-Process -NoNewWindow -FilePath "powershell.exe" -ArgumentList "go run .\vehicle_registration_service.go"
cd ..

## Running of Stripe Service
Set-Location .\stripe-service\
Write-Host "Starting Vehicle Service..."
Start-Process -NoNewWindow -FilePath "powershell.exe" -ArgumentList "go run .\stripe_service.go"
cd ..

# Go Back to CNAD_Assg1
cd ..

Write-Host "All Microservices are up and running!!!"

exit
