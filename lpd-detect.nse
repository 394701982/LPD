-- lpd-check.nse
--
-- This script detects a custom LPD (Line Printer Daemon) protocol implementation.
-- It sends a specific command and checks if the response matches the expected pattern.

local nmap = require "nmap"
local shortport = require "shortport"
local stdnse = require "stdnse"
local string = require "string"

description = [[
  Detects a custom LPD protocol service by sending a specific command and checking the response.
]]

author = "Your Name"
license = "Same as Nmap--See https://nmap.org/book/man-legal.html"
categories = {"discovery", "safe"}

portrule = shortport.port_or_service(515, "printer")

action = function(host, port)
  local socket = nmap.new_socket()
  socket:set_timeout(5000)

  local status, err = socket:connect(host, port)
  if not status then
    return stdnse.format_output(false, "Connection failed: %s", err)
  end

  -- Send the command
  local command = string.char(0x02) .. "print_queue\n"
  socket:send(command)

  -- Receive the response
  local response, err = socket:receive_lines(1)
  socket:close()

  if not response then
    return stdnse.format_output(false, "No response from the server: %s", err)
  end

  if response then
    return ("Detected custom LPD service: %s"):format(response)
  else
    return "Tor Control Service detected but no response received."
  end
end
