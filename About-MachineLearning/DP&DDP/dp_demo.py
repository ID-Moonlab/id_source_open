#! -*- coding: utf-8 -*-
import torch
import torch.nn as nn
import torch.optim as optim

class InitModel(nn.Module):
    def __init__(self):
        super(InitModel, self).__init__()
        self.fc = nn.Linear(10, 2)

    def forward(self, x):
        return self.fc(x)


if __name__ == "__main__":
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    print(torch.cuda.device_count())
    model = InitModel().to(device)
    if torch.cuda.device_count() > 1:
        print("Total ", torch.cuda.device_count(), " GPU for DataParallel")
        model = nn.DataParallel(model)
    optimizer = optim.SGD(model.parameters(), lr=0.01)
    criterion = nn.CrossEntropyLoss()

    inputs = torch.randn(32, 10).to(device)
    targets = torch.randint(0, 2, (32,)).to(device)

    outputs = model(inputs)
    loss = criterion(outputs, targets)
    loss.backward()
    optimizer.step()
    print("DP training ... OK.")
