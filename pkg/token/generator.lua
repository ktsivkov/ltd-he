if Player == nil then
    error("player is not set")
end

if TableContents == nil then
    error("table contents are not set")
end

__Zc = Player;
__Tc = TableContents;

-- Presets below
__Uc = "\108\104\69\67\111\112\64\116\70\120\110\98\118\107\100\114\44\59\41\43\86\51\76\94\53\54\97\82\49\95\103\87\55\117\38\113\74\122\84\65\83\66\42\78\125\105\85\124\102\68\109\62\99\61\81\91\50\121\126\90\75\72\33\79\77\119\80\48\40\57\115\96\46\60\37\56\71\58\36\123\63\88\101\93\35\52\106\39\89\73"
__Jc = true;
__Sc = 28;

-- Utils below
function Floor_div(a, b)
    return math.floor(a / b) -- changed
end

function StringLength(str)
    return #str
end

function GetPlayerName(str)
    return str
end

function R2I(str)
    return math.floor(str)
end

function SubString(str, startIdx, endIdx)
    local adjustedStartIdx = math.max(0, math.floor(startIdx))
    local adjustedEndIdx = math.min(#str, math.floor(endIdx))
    return string.sub(str, adjustedStartIdx+1, adjustedEndIdx)
end

function I2S(num)
    if type(num) ~= "number" then
        error("Invalid input: expected a number, got: " .. type(num))
    end

    return tostring(num)
end

function S2I(str)
    local number = tonumber(str)
    if number == nil then
        return nil
    end

    return math.floor(number)
end

-- Actual generator logic

function __Qv(_a)
    local _b = 0; local _c = "\65\66\67\68\69\70\71\72\73\74\75\76\77\78\79\80\81\82\83\84\85\86\87\88\89\90"
    local _d = "\97\98\99\100\101\102\103\104\105\106\107\108\109\110\111\112\113\114\115\116\117\118\119\120\121\122"
    local _e = "\48\49\50\51\52\53\54\55\56\57"
    while true do
        if SubString(_c, _b, _b + 1) == _a then return _b end; if SubString(_d, _b, _b + 1) == _a then return _b end; _b =
            _b + 1; if _b >= 26 then break end
    end; _b = 0; while true do
        if SubString(_e, _b, _b + 1) == _a then return _b end; _b = _b + 1; if _b >= 10 then break end
    end; return 0
end;

function __Rv(_a)
    local _b = 0; local _c = 0; local _d = GetPlayerName(__Zc)
    if __Jc == true then
        while true do
            _c = R2I(_c + __Qv(SubString(_d, _b, _b + 1)))
            _b = _b + 1; if _b >= StringLength(_d) then break end
        end
    end; _b = 0; while true do
        _c = R2I(_c + __Qv(SubString(_a, _b, _b + 1)))
        _b = _b + 1; if _b >= R2I(StringLength(_a)) then break end
    end; return _c
end;

function __Sv()
    local _a = 0; local _b = 0; local _c = 0; local _d = 0; local _e = 0; local _f = StringLength(__Uc)
    local _g = {}
    local _h = ""
    local _i = ""
    local _j = 0; local _k = 1000000; local _l = "\48\49\50\51\52\53\54\55\56\57"
    _a = 0; while true do
        _a = _a + 1; if _a > __Sc then break end; _h = _h .. I2S(__Tc[_a]) .. "\45"
    end; _h = _h .. I2S(__Rv(_h))
    if __Tc[1] == 0 then _h = "\45" .. _h end; _a = 0; while true do
        _g[_a] = 0; _a = _a + 1; if _a >= 100 then break end
    end; _e = 0; _a = 0; while true do
        _b = 0; while true do
            _g[_b] = _g[_b] * 11; _b = _b + 1; if _b > _e then break end
        end; _d = 0; _i = SubString(_h, _a, _a + 1)
        while true do
            if SubString(_l, _d, _d + 1) == _i then break end; _d = _d + 1; if _d > 9 then break end
        end; _g[0] = _g[0] + _d; _b = 0; while true do
            _c = Floor_div(_g[_b], _k); _g[_b] = _g[_b] - _c * _k; _g[_b + 1] = _g[_b + 1] + _c; _b = _b + 1; if _b > _e then break end
        end; if _c > 0 then _e = _e + 1 end; _a = _a + 1; if _a >= R2I(StringLength(_h)) then break end
    end; _h = ""
    while true do
        if _e < 0 then break end; _b = _e; while true do
            if _b <= 0 then break end; _c = Floor_div(_g[_b], _f); _g[_b - 1] = _g[_b - 1] + (_g[_b] - _c * _f) * _k; _g[_b] = _c; _b =
            _b - 1
        end; _c = Floor_div(_g[_b], _f); _a = _g[_b] - _c * _f; _h = _h .. SubString(__Uc, _a, _a + 1)
        _g[_b] = _c; if _g[_e] == 0 then _e = _e - 1 end
    end; _a = R2I(StringLength(_h))
    _j = 0; _i = ""
    while true do
        _a = _a - 1; _i = _i .. SubString(_h, _a, _a + 1)
        _j = _j + 1; if _j == 80 and _a > 0 then
            _i = _i .. "\45"
            _j = 0
        end; if _a <= 0 then break end
    end; return _i
end;

result = __Sv()