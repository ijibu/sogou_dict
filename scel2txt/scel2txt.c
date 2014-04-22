#include <stdio.h>
#include <string.h>
#include <malloc.h>
#include <memory.h>
/*
linux下批量转换scel文件到txt的方式
cd scel目录
find ./ -name "*scel" -exec ./scel2txt {} \ ;
*/

typedef struct PY_
{
        unsigned short mark;
        char py[6+1];
        struct PY_ * next;
} PY;

int unicode2utf8char(unsigned short in, unsigned char * out)
{
	//out 为unsigned char str[16]={0};		//刚好两个字节
    if (in >= 0x0000 && in <= 0x007f)
    {
        *out=in;
        return 0;
    }
    else if (in >= 0x0080 && in <= 0x07ff)
    {
        *out = 0xc0 | (in >> 6);
        out ++;
        *out = 0x80 | (in & (0xff >> 2));
        return 0;
    }
    else if (in >= 0x0800 && in <= 0xffff)
    {
        *out = 0xe0 | (in >> 12);
        out ++;
        *out = 0x80 | (in >> 6 & 0x003f);
        out ++;
        *out = 0x80 | (in & (0xff >> 2));
        return 0;
    }
    printf("输入的不是short吧,解析有问题\n");
    return 0;
}

int unicode2utf8str(char * in, int insize,unsigned char * out)
{
    unsigned char str[16]={0};		//刚好两个字节
    unsigned short tmp[insize/2];
    int i;

    *out='\0';
    memcpy(tmp,in,insize);	//拷贝in所指向的内存前insize个字节tmp

    for( i=0;i<insize/2;i++)
    {
        memset(str,0,sizeof(str));	//内存初始化，两个字节全部设置为0。
		printf("%X", tmp[i]);		//5BA0
        printf("\n");
        unicode2utf8char(tmp[i],str);
        printf("%X", str);			//351EB9D0
        printf("\n");
        strcat(out,str);	//连接两字符串
    }
    return 0;
}


PY * loadPY(FILE * fp)
{
    unsigned char str[128]={0};
    unsigned char outstr[128]={0};
    unsigned short num[16]={0};
    int i;
    PY * head=NULL;
    PY * p=NULL;

    fseek(fp, 0x1540, SEEK_SET);
    fgets(str,4+1,fp);

    if(memcmp(str,"\x9D\x01\x00\x00",4) != 0)
    {
            printf("莫非解析位置有误?\n");
            //return -1;

    }

    head=(PY *)malloc(sizeof(PY));
    head->next=NULL;
    p=head;
    while(1)
    {
            memset(str,0,sizeof(str));
            memset(num,0,sizeof(num));
            for(i=0;i<4;i++)
            {
                    str[i]=fgetc(fp);
            }
            memcpy(num,str,4);

            p->next=(PY *)malloc(sizeof(PY));
            p=p->next;
            p->mark=num[0];

            memset(str,0,sizeof(str));
            fgets(str,num[1]+1,fp);
            unicode2utf8str(str,64,p->py);

            p->next=NULL;
            if( strcmp(p->py,"zuo" ) == 0)
            {
                    return head;
                    break;
            }
    }
}

int creatWordStock(FILE *fp,PY * head, char * fileName)
{
    unsigned char str[256]={0};
    unsigned char outstr[256]={0};
    unsigned char pybuf[128]={0};
    unsigned char hzbuf[128]={0};
    unsigned char buf[256]={0};
    PY *p =NULL;
    FILE * newfp;
    unsigned short num[64]={0};
    int i,count,offset;

    newfp=fopen(fileName,"w+");
    if( newfp == NULL)
    {
            perror("fopen error");
            return -1;
    }

    fseek(fp, 0x2628, SEEK_SET);
    while(1)
    {
            count=0;
            offset=0;
            p=head->next;
            memset(num,0,sizeof(num));
            memset(str,0,sizeof(str));
            memset(pybuf,0,sizeof(pybuf));
            memset(hzbuf,0,sizeof(hzbuf));
            memset(buf,0,sizeof(buf));

            for(i=0;i<4;i++)
            {
                    str[i]=fgetc(fp);
                    if( feof(fp) )
                    {
                            fclose(newfp);
                            return 0;
                    }
            }

            memcpy(num,str,4);
            offset=num[0]-1;
            count=num[1];
            memset(str,0,sizeof(str));
            for(i=0;i<count;i++)
            {
                    str[i]=fgetc(fp);
                    if( feof(fp) )
                    {
                            fclose(newfp);
                            return 0;
                    }
            }
            memset(num,0,sizeof(num));
            memcpy(num,str,count);

            for(i=0;i<count/2;i++)
            {
                    p=head->next;
                    while(p!=NULL)
                    {
                            if( p->mark == num[i])
                            {
                                    strcat(pybuf,p->py);
                                    strcat(pybuf,"'");
                                    p=NULL;
                                    break;
                            }
                            p=p->next;
                    }
            }
            if( pybuf[strlen(pybuf)-1] == '\'' )
                    pybuf[strlen(pybuf)-1] = '\0';

            memset(num,0,sizeof(num));
            memcpy(num,str,count);
            for(i=0;i<2;i++)
            {
                    str[i]=fgetc(fp);
                    if( feof(fp) )
                    {
                            fclose(newfp);
                            return 0;
                    }
            }
            memcpy(num,str,count);
            count=num[0];

            memset(num,0,sizeof(num));
            memcpy(num,str,count);
            for(i=0;i<count;i++)
            {
                    str[i]=fgetc(fp);
                    if( feof(fp) )
                    {
                            fclose(newfp);
                            return 0;
                    }
            }
            unicode2utf8str(str,64,hzbuf);
            sprintf(buf,"%s %s",pybuf,hzbuf);
            fprintf(newfp,"%s\n",buf);
            for(i=0;i<(12+offset*(12+count+2));i++)
            {
                    str[i]=fgetc(fp);
                    if( feof(fp) )
                    {
                            fclose(newfp);
                            return 0;
                    }
            }
    }
    return 0;
}

void freePY(PY * head)
{
    PY * p;
    p=head;
    if( p->next !=NULL)
    {
        head=p;
        p=p->next;
        free(head);
    }
}


int main(int argc ,char * argv[])
{
    FILE * fp=NULL;
    unsigned char str[128]={0};
    unsigned char outstr[128]={0};
    PY * head;
    int i;

    if(argc <=1)
    {
        printf("请输入sg词库文件!");
        return 0;
    }
    fp=fopen(argv[1],"r");
    if( fp == NULL)
    {
        perror("fopen error");
        return -1;
    }
    fgets(str,8+1,fp);
    if( memcmp(str,"\x40\x15\x00\x00\x44\x43\x53\x01",8))	//比较内存区域
    {
        printf("你确认你选择的是搜狗(.scel)词库?\n");
        return 0;
    }
    memset(str,0,sizeof(str));	//将str指向的内存空间前sizeof(str)个字节填入值0

    fseek(fp, 0x130, SEEK_SET);
    fgets(str,64+1,fp);
    unicode2utf8str(str,64,outstr);
    printf("字库名称:%s\n",outstr);

    memset(str,0,sizeof(str));
    memset(outstr,0,sizeof(outstr));
    fseek(fp, 0x338, SEEK_SET);
    fgets(str,64+1,fp);
    unicode2utf8str(str,64,outstr);
    printf("字库类别:%s\n",outstr);

    memset(str,0,sizeof(str));
    memset(outstr,0,sizeof(outstr));
    fseek(fp, 0x540, SEEK_SET);
    fgets(str,64+1,fp);
    unicode2utf8str(str,64,outstr);
    printf("字库信息:%s\n",outstr);

    memset(str,0,sizeof(str));
    memset(outstr,0,sizeof(outstr));
    fseek(fp, 0xd40, SEEK_SET);
    fgets(str,64+1,fp);

    unicode2utf8str(str,64,outstr);
    printf("字库示例:%s\n",outstr);

    head=loadPY(fp);
    char * fileName = strcat(argv[1],".txt");
    creatWordStock(fp,head,fileName);
    freePY(head);

    fclose(fp);
    return 0;
}