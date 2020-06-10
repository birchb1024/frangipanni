-- Beginning of XML output.
-- Unlikely for input data to be valid in XML though.
-- Example:
-- <root count="1" sep="">
--    <usr/local/go count="7" sep="/">
--       <test/fixedbugs/bug028.go count="1" sep="/"/>
--       <test2 count="4" sep="/">
--          <fixedbugs1 count="2" sep="/">
--             <bug205.go count="1" sep="/"/>
--             <bug206.go count="1" sep="/"/>


function tablelength(T)
    local count = 0
    for _ in pairs(T) do count = count + 1 end
    return count
end

function indent(n)
    for i=1, n do
        io.write("   ")
    end
end



function xml(node)
    indent(node.depth)
    io.write("<")
    io.write(node.text)
    io.write(" count=\""..node.numMatched.."\"")
    io.write(" sep=\""..node.sep.."\"")
    if tablelength(node.children) == 0 then
        io.write("/>\n")
        return
    end
    io.write(">\n")
    local tkeys = {}
    -- populate the table that holds the keys
    for k in pairs(node.children) do table.insert(tkeys, k) end
    -- sort the keys
    table.sort(tkeys)
    -- use the keys to retrieve the values in the sorted order
    for _, k in ipairs(tkeys) do 
        xml(node.children[k]) 
    end
    indent(node.depth)
    io.write("</")
    io.write(node.text)
    io.write(">\n")
end
frangipanni.text = "root"
xml(frangipanni)
