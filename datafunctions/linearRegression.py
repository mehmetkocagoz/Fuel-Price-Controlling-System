import pandas as pd
import matplotlib.pyplot as plt
from sklearn.model_selection import train_test_split
from sklearn.linear_model import LinearRegression
from sklearn import metrics
import numpy as np


dataFrame2 = pd.read_csv("prices.csv")
print(dataFrame2.corr(method="pearson",numeric_only=True))

dataFrame2['brentoilpriceTL'] = dataFrame2['brentoilprice'] * dataFrame2['exchangerate']

# Save the modified DataFrame to a new CSV file
dataFrame2.to_csv('output_file.csv', header=True, index=False)

x2 = dataFrame2[['brentoilprice','exchangerate']]
y2 = dataFrame2['fuelprice']
X2_train, X2_test, y2_train, y2_test = train_test_split(x2, y2, test_size=0.2, random_state=101)
lm2 = LinearRegression()
lm2.fit(X2_train,y2_train)
coeff_df = pd.DataFrame(lm2.coef_,x2.columns,columns=['Coefficient'])
print(coeff_df)
predictions2 = lm2.predict(X2_test)
#plt.scatter(y2_test,predictions)
#plt.show()
MAE2 = metrics.mean_absolute_error(y2_test,predictions2)
MSE2 = metrics.mean_squared_error(y2_test,predictions2)
RMSE2 = np.sqrt(metrics.mean_squared_error(y2_test,predictions2)) 
print('MAE:',MAE2)
print('MSE:',MSE2)
print('RMSE:',RMSE2)

dataFrame = pd.read_csv("output_file.csv")
# Set the maximum display rows and columns to a high value
pd.set_option('display.max_rows', 1000)
pd.set_option('display.max_columns', 1000)
print(dataFrame.corr(method="pearson",numeric_only=True))

x = dataFrame[['brentoilpriceTL','exchangerate']]
y = dataFrame[['fuelprice']]
X_train, X_test, y_train, y_test = train_test_split(x, y, test_size=0.2, random_state=101)
lm = LinearRegression()
lm.fit(X_train,y_train)
predictions = lm.predict(X_test)
#plt.scatter(y2_test,predictions)
#plt.show()
MAE = metrics.mean_absolute_error(y_test,predictions)
MSE = metrics.mean_squared_error(y_test,predictions)
RMSE = np.sqrt(metrics.mean_squared_error(y_test,predictions)) 
print('MAE:',MAE)
print('MSE:',MSE)
print('RMSE:',RMSE)


'''
By doing same operations on clean data, we can see that the results are better.
MAE, MSE and RMSE values are lower than the previous ones.
It means that the model is more accurate.
'''
