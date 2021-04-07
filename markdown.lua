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

function markdown(node, bullet)
    if node.lineNumber > 0 then  -- don't write a root note
        indent(node.depth -1)
        io.write(bullet)
        print(node.text)
    end
    for k, v in pairs(node.children) do
        markdown(v, bullet)
    end
end

markdown(frangipanni, os.getenv("NUMBERED_LIST") and "1. " or "* ")
