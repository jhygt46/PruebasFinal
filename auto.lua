request = function()
    path = '/a?p=1&s=aa'
    return wrk.format("GET", path)
end