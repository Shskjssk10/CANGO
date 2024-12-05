# Print a message before restarting the terminal
Write-Host "Restarting the terminal..."

# Terminate the current PowerShell session
Stop-Process -Id $PID

exit
