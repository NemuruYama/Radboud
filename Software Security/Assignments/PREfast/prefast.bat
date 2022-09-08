@echo off

IF "%*" == "" goto noparam
sdkvcbin\cl.exe /analyze /W 3 /I sdkinclude /I sdkvcinclude /c %*
goto end

:noparam
sdkvcbin\cl.exe /analyze /W 3 /I sdkinclude /I sdkvcinclude /c prefast_exercise.cpp

:end
