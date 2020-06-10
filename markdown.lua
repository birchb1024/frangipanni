-- Beginning of Markdown unnumbered indented bullet list
-- Print each line preceded with indents and an '*'
-- Example:
-- * 
--    * x.a
--       * 1
--       * 2
--    * Z
--    * 1.2
--    * A
--    * C
--       * D
--       * 2

function indent(n)
    for i=1, n do
        io.write("   ")
    end
end

function markdown(node)
    indent(node.depth)
    io.write("* ")
    print(node.text)
    for k, v in pairs(node.children) do
        markdown(v)
    end
end

markdown(frangipanni)
