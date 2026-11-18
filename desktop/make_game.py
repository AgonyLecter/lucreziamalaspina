import sys,shutil,os
src=os.path.dirname(os.path.dirname(os.path.abspath(__file__))); dst=sys.argv[1]; with_sw=len(sys.argv)>2
shutil.rmtree(dst,ignore_errors=True); os.makedirs(dst)
s=open(src+'/index.html',encoding='utf-8').read()
old_link=[l for l in s.split('\n') if 'fonts.googleapis.com/css2' in l][0]
s=s.replace(old_link,'<link href="fonts/fonts.css" rel="stylesheet">')
s=s.replace('<link rel="preconnect" href="https://fonts.googleapis.com">\n','').replace('<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>\n','')
open(dst+'/index.html','w',encoding='utf-8').write(s)
for f in ['manifest.webmanifest','icon-192.png','icon-512.png','apple-touch-icon.png','favicon.png']+(['sw.js'] if with_sw else []):
    shutil.copy(src+'/'+f,dst+'/'+f)
shutil.copytree(os.path.dirname(os.path.abspath(__file__))+'/fonts',dst+'/fonts')
print('ok',dst)
