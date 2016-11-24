# Snowden
This little tool allows users to be notified when one or more of their watched files or folders are included in a new or updated 
Github pull request, so they can review the changes made on them.

## Dependencies

Go > 1.5

This program is designed to work in conjunction with https://github.com/adnanh/webhook. Take a look at its documentation to know
more about how install and configure it.

## Installation

`go get github.com/svera/snowden`

## Usage

Snowden expects 5 parameters, in this order:

* 