if Player == nil then
    error("player is not set")
end

if Token == nil then
    error("token is not set")
end

__Zc = Player;
-- Presets below
__Tc = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
__Uc = "lhECop@tFxnbvkdr,;)+V3L^56aR1_gW7u&qJzTASB*N}iU|fDm>c=Q[2y~ZKH!OMwP0(9s`.<%8G:${?Xe]#4j'YI"
__Jc = true;
__Sc = 28

-- Utils below
function Floor_div(a, b)
    return a / b -- changed
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

-- Actual validator logic

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

function __Tv(_a)
    local _b; local _c; local _d; local _e; local _f = 0; local _g; local _h = {}
    local _i = ""
    local _j = StringLength(__Uc)
    local _k = -1; local _l = 1000000; local _m = "0123456789-"
    local _n; _b = 0; while true do
        _h[_b] = 0; _b = _b + 1; if _b >= 100 then break end
    end; _g = 0; _b = 0; while true do
        _c = 0; while true do
            _h[_c] = _h[_c] * _j; _c = _c + 1; if _c > _g then break end
        end; _k = _k + 1; if _k == 80 then
            _k = 0; _b = _b + 1
        end; _e = _j; _n = SubString(_a, _b, _b + 1)
        while true do
            _e = _e - 1; if _e < 1 then break end; if SubString(__Uc, _e, _e + 1) == _n then break end
        end; _h[0] = _h[0] + _e; _c = 0; while true do
            _d = R2I(Floor_div(_h[_c], _l))
            _h[_c] = _h[_c] - _d * _l; _h[_c + 1] = _h[_c + 1] + _d; _c = _c + 1; if _c > _g then break end
        end; if _d > 0 then _g = _g + 1 end; _b = _b + 1; if _b >= StringLength(_a) then break end
    end; while true do
        if _g < 0 then break end; _c = _g; while true do
            if _c <= 0 then break end; _d = R2I(Floor_div(_h[_c], 11))
            _h[_c - 1] = _h[_c - 1] + (_h[_c] - _d * 11) * _l; _h[_c] = _d; _c = _c - 1
        end; _d = R2I(Floor_div(_h[_c], 11))
        _b = _h[_c] - _d * 11; _i = SubString(_m, _b, _b + 1) .. _i; _h[_c] = _d; if _h[_g] == 0 then _g = _g - 1 end
    end; _b = 0; _c = 0;

    while true do
        while true do
            if _b >= R2I(StringLength(_i)) then break end;
            if _b > 0 and SubString(_i, _b, _b + 1) == "\45" and SubString(_i, _b - 1, _b) ~= "\45" then break end;
            _b = _b + 1
        end;
        if _b < R2I(StringLength(_i)) then _d = _b end; _f = _f + 1; __Tc[_f] = S2I(SubString(_i, _c, _b))
        _c = _b + 1; _b = _b + 1; if _b >= R2I(StringLength(_i)) then break end
    end; _c = R2I(__Rv(SubString(_i, 0, _d)))

    __Sc = _f - 1;
    if _c == __Tc[_f] then return true end; return false
end;

result = __Tv(Token)

--error(__Tc[6] .. __Tc[2])
