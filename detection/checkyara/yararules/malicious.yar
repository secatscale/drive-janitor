rule MaliciousFile
{
    strings:
        $mal = "verymalicious"
    condition:
        $mal
}
