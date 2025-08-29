// todo: replace your os image address
FROM xxx

USER work

# 复制output下所有文件到工作目录，不包括output目录本身
COPY --chown=work output/ /home/work/mcp-host-demo/

WORKDIR /home/work/mcp-host-demo

EXPOSE 8083

CMD ["bin/mcp-host-demo"]