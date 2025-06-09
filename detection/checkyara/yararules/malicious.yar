rule MaliciousFile
{
    strings:
        $mal = "malicious"
    condition:
        $mal
}