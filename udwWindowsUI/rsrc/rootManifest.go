package rsrc

var gRootManifestContent = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
    <assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0" xmlns:asmv3="urn:schemas-microsoft-com:asm.v3">
    	<assemblyIdentity version="9.0.0.0" processorArchitecture="*" name="foxWinUI.exe" type="win32"/>
        <description>A tool</description>
        <dependency>
            <dependentAssembly>
                <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
            </dependentAssembly>
        </dependency>
    	<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
            <security>
                <requestedPrivileges>
                    <requestedExecutionLevel level="requireAdministrator"/>
                </requestedPrivileges>
            </security>
        </trustInfo>
    </assembly>`)

var gNormalManifestContent = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
    <assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0" xmlns:asmv3="urn:schemas-microsoft-com:asm.v3">
    	<assemblyIdentity version="9.0.0.0" processorArchitecture="*" name="foxWinUI.exe" type="win32"/>
        <description>A tool</description>
        <dependency>
            <dependentAssembly>
                <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
            </dependentAssembly>
        </dependency>
    	<asmv3:application>
    		<asmv3:windowsSettings xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">
    			<dpiAware>true</dpiAware>
    		</asmv3:windowsSettings>
    	</asmv3:application>
    </assembly>`)
